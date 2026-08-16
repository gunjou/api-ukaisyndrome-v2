package jadwal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) GetJadwalPeserta(ctx context.Context, userID int) ([]JadwalDTO, error) {
	query := `
		SELECT
			j.id_jadwal, j.id_paketkelas, pk.nama_kelas,
			j.id_mentor, u.nama AS nama_mentor, u.nickname AS nickname_mentor,

			TO_CHAR(j.tanggal, 'YYYY-MM-DD') AS tanggal,
			TO_CHAR(j.waktu_mulai, 'HH24:MI:SS') AS waktu_mulai,
			TO_CHAR(j.waktu_selesai, 'HH24:MI:SS') AS waktu_selesai,

			TO_CHAR(j.tanggal_reschedule, 'YYYY-MM-DD') AS tanggal_reschedule,
			TO_CHAR(j.waktu_mulai_reschedule, 'HH24:MI:SS') AS waktu_mulai_reschedule,
			TO_CHAR(j.waktu_selesai_reschedule, 'HH24:MI:SS') AS waktu_selesai_reschedule,

			TO_CHAR(COALESCE(j.tanggal_reschedule, j.tanggal), 'YYYY-MM-DD') AS tanggal_efektif,
			TO_CHAR(COALESCE(j.waktu_mulai_reschedule, j.waktu_mulai), 'HH24:MI:SS') AS waktu_mulai_efektif,
			TO_CHAR(COALESCE(j.waktu_selesai_reschedule, j.waktu_selesai), 'HH24:MI:SS') AS waktu_selesai_efektif,

			j.type_pertemuan
		FROM jadwal_kelas j
		INNER JOIN paketkelas pk
			ON pk.id_paketkelas = j.id_paketkelas
		INNER JOIN pesertakelas pks
			ON pks.id_paketkelas = j.id_paketkelas
		INNER JOIN users u
			ON u.id_user = j.id_mentor
		WHERE pks.id_user = $1
		  AND j.status = 1
		  AND pk.status = 1
		  AND pks.status = 1
		  AND u.status = 1
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

	var result []JadwalDTO

	for rows.Next() {
		var j JadwalDTO

		if err := rows.Scan(
			&j.IDJadwal, &j.IDPaketKelas, &j.NamaKelas,
			&j.IDMentor, &j.NamaMentor, &j.NicknameMentor,
			&j.Tanggal, &j.WaktuMulai, &j.WaktuSelesai,
			&j.TanggalReschedule, &j.WaktuMulaiReschedule, &j.WaktuSelesaiReschedule,
			&j.TanggalEfektif, &j.WaktuMulaiEfektif, &j.WaktuSelesaiEfektif,
			&j.TypePertemuan,
		); err != nil {
			return nil, err
		}

		result = append(result, j)
	}

	return result, rows.Err()
}


func (r *Repository) GetJadwalPesertaByID(ctx context.Context, userID, jadwalID int) (*JadwalDTO, error) {
	query := `
		SELECT
			j.id_jadwal, j.id_paketkelas, pk.nama_kelas,
			j.id_mentor, u.nama AS nama_mentor, u.nickname AS nickname_mentor,

			TO_CHAR(j.tanggal, 'YYYY-MM-DD') AS tanggal,
			TO_CHAR(j.waktu_mulai, 'HH24:MI:SS') AS waktu_mulai,
			TO_CHAR(j.waktu_selesai, 'HH24:MI:SS') AS waktu_selesai,

			TO_CHAR(j.tanggal_reschedule, 'YYYY-MM-DD') AS tanggal_reschedule,
			TO_CHAR(j.waktu_mulai_reschedule, 'HH24:MI:SS') AS waktu_mulai_reschedule,
			TO_CHAR(j.waktu_selesai_reschedule, 'HH24:MI:SS') AS waktu_selesai_reschedule,

			TO_CHAR(COALESCE(j.tanggal_reschedule, j.tanggal), 'YYYY-MM-DD') AS tanggal_efektif,
			TO_CHAR(COALESCE(j.waktu_mulai_reschedule, j.waktu_mulai), 'HH24:MI:SS') AS waktu_mulai_efektif,
			TO_CHAR(COALESCE(j.waktu_selesai_reschedule, j.waktu_selesai), 'HH24:MI:SS') AS waktu_selesai_efektif,

			j.type_pertemuan
		FROM jadwal_kelas j
		INNER JOIN paketkelas pk
			ON pk.id_paketkelas = j.id_paketkelas
		INNER JOIN pesertakelas pks
			ON pks.id_paketkelas = j.id_paketkelas
		INNER JOIN users u
			ON u.id_user = j.id_mentor
		WHERE j.id_jadwal = $1
		  AND pks.id_user = $2
		  AND j.status = 1
		  AND pk.status = 1
		  AND pks.status = 1
		  AND u.status = 1
	`

	var j JadwalDTO

	err := r.DB.QueryRow(ctx, query, jadwalID, userID).Scan(
		&j.IDJadwal, &j.IDPaketKelas, &j.NamaKelas,
		&j.IDMentor, &j.NamaMentor, &j.NicknameMentor,
		&j.Tanggal, &j.WaktuMulai, &j.WaktuSelesai,
		&j.TanggalReschedule, &j.WaktuMulaiReschedule, &j.WaktuSelesaiReschedule,
		&j.TanggalEfektif, &j.WaktuMulaiEfektif, &j.WaktuSelesaiEfektif,
		&j.TypePertemuan,
	)

	if err != nil {
		return nil, err
	}

	return &j, nil
}