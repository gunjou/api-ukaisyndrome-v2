package jadwal

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrJadwalNotFound = errors.New("jadwal tidak ditemukan")

type Repository struct {
	DB *pgxpool.Pool
}

// baseJadwalQuery: $1 selalu userID (peserta).
//
// Aturan akses:
//   - ONLINE  -> peserta terdaftar di paket kelas jadwal (pesertakelas)
//   - OFFLINE -> peserta anggota aktif kelompok belajar jadwal (kelompok_belajar_peserta)
//
// Pakai EXISTS (bukan INNER JOIN ke tabel anggota) supaya tidak ada baris duplikat
// dan supaya kedua jalur akses bisa digabung dengan OR.
const baseJadwalQuery = `
	SELECT
		j.id_jadwal,

		j.id_paketkelas,
		pk.nama_kelas,
		j.id_kelompok,
		kb.nama_kelompok,

		j.id_mentor,
		u.nama AS nama_mentor,
		u.nickname AS nickname_mentor,

		j.topik,
		j.catatan,

		TO_CHAR(j.tanggal, 'YYYY-MM-DD') AS tanggal,
		TO_CHAR(j.waktu_mulai, 'HH24:MI:SS') AS waktu_mulai,
		TO_CHAR(j.waktu_selesai, 'HH24:MI:SS') AS waktu_selesai,

		TO_CHAR(j.tanggal_reschedule, 'YYYY-MM-DD') AS tanggal_reschedule,
		TO_CHAR(j.waktu_mulai_reschedule, 'HH24:MI:SS') AS waktu_mulai_reschedule,
		TO_CHAR(j.waktu_selesai_reschedule, 'HH24:MI:SS') AS waktu_selesai_reschedule,

		TO_CHAR(COALESCE(j.tanggal_reschedule, j.tanggal), 'YYYY-MM-DD') AS tanggal_efektif,
		TO_CHAR(COALESCE(j.waktu_mulai_reschedule, j.waktu_mulai), 'HH24:MI:SS') AS waktu_mulai_efektif,
		TO_CHAR(COALESCE(j.waktu_selesai_reschedule, j.waktu_selesai), 'HH24:MI:SS') AS waktu_selesai_efektif,

		j.type_pertemuan,

		ap.id_absensi_peserta,
		ap.status_kehadiran,
		TO_CHAR(ap.check_in_at, 'YYYY-MM-DD HH24:MI:SS') AS check_in_at,
		ap.jenis_absensi::text,
		ap.status_verifikasi::text,
		ap.bukti_kehadiran_url,
		(ap.id_absensi_peserta IS NOT NULL) AS sudah_absen

	FROM jadwal_kelas j

	LEFT JOIN paketkelas pk
		ON pk.id_paketkelas = j.id_paketkelas
		AND pk.status = 1

	LEFT JOIN kelompok_belajar kb
		ON kb.id_kelompok = j.id_kelompok
		AND kb.status = 1

	INNER JOIN users u
		ON u.id_user = j.id_mentor
		AND u.status = 1

	LEFT JOIN absensi_peserta ap
		ON ap.id_jadwal = j.id_jadwal
		AND ap.id_peserta = $1
		AND ap.status = 1

	WHERE j.status = 1
	  AND (
		(
			UPPER(j.type_pertemuan) = 'ONLINE'
			AND pk.id_paketkelas IS NOT NULL
			AND EXISTS (
				SELECT 1
				FROM pesertakelas pks
				WHERE pks.id_paketkelas = j.id_paketkelas
				  AND pks.id_user = $1
				  AND pks.status = 1
			)
		)
		OR
		(
			UPPER(j.type_pertemuan) = 'OFFLINE'
			AND kb.id_kelompok IS NOT NULL
			AND EXISTS (
				SELECT 1
				FROM kelompok_belajar_peserta kbp
				WHERE kbp.id_kelompok = j.id_kelompok
				  AND kbp.id_peserta = $1
				  AND kbp.status = 1
			)
		)
	  )
`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanJadwal(row rowScanner, j *JadwalDTO) error {
	return row.Scan(
		&j.IDJadwal,

		&j.IDPaketKelas,
		&j.NamaKelas,
		&j.IDKelompok,
		&j.NamaKelompok,

		&j.IDMentor,
		&j.NamaMentor,
		&j.NicknameMentor,

		&j.Topik,
		&j.Catatan,

		&j.Tanggal,
		&j.WaktuMulai,
		&j.WaktuSelesai,

		&j.TanggalReschedule,
		&j.WaktuMulaiReschedule,
		&j.WaktuSelesaiReschedule,

		&j.TanggalEfektif,
		&j.WaktuMulaiEfektif,
		&j.WaktuSelesaiEfektif,

		&j.TypePertemuan,

		&j.IDAbsensiPeserta,
		&j.StatusKehadiran,
		&j.CheckInAt,
		&j.JenisAbsensi,
		&j.StatusVerifikasi,
		&j.BuktiKehadiranURL,
		&j.SudahAbsen,
	)
}

func (r *Repository) GetJadwalPeserta(ctx context.Context, userID int) ([]JadwalDTO, error) {
	query := baseJadwalQuery + `
		ORDER BY
			COALESCE(j.tanggal_reschedule, j.tanggal) ASC,
			COALESCE(j.waktu_mulai_reschedule, j.waktu_mulai) ASC,
			j.id_jadwal ASC
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// make(..., 0) supaya response JSON berupa [] bukan null saat kosong
	result := make([]JadwalDTO, 0)

	for rows.Next() {
		var j JadwalDTO
		if err := scanJadwal(rows, &j); err != nil {
			return nil, err
		}
		result = append(result, j)
	}

	return result, rows.Err()
}

func (r *Repository) GetJadwalPesertaByID(ctx context.Context, userID, jadwalID int) (*JadwalDTO, error) {
	query := baseJadwalQuery + `
		AND j.id_jadwal = $2
	`

	var j JadwalDTO
	if err := scanJadwal(r.DB.QueryRow(ctx, query, userID, jadwalID), &j); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJadwalNotFound
		}
		return nil, err
	}

	return &j, nil
}