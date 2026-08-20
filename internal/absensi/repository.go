package absensi

import (
	"api-ukaisyndrome-v2/pkg/timeutil"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}


func (r *Repository) IsJadwalPeserta(
	ctx context.Context,
	userID, jadwalID int,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM jadwal_kelas j
			INNER JOIN pesertakelas pks
				ON pks.id_paketkelas = j.id_paketkelas
			WHERE j.id_jadwal = $1
			  AND pks.id_user = $2
			  AND j.status = 1
			  AND pks.status = 1
		)
	`

	var exists bool
	err := r.DB.QueryRow(ctx, query, jadwalID, userID).Scan(&exists)

	return exists, err
}


func (r *Repository) GetAbsensiByJadwal(
	ctx context.Context,
	userID, jadwalID int,
) (*AbsensiPesertaDTO, error) {
	query := `
		SELECT
			ap.id_absensi_peserta, ap.id_jadwal,
			j.id_paketkelas, pk.nama_kelas, j.topik, j.catatan,
			ap.id_peserta, u.nama, u.nickname,
			ap.status_kehadiran,
			TO_CHAR(ap.check_in_at, 'YYYY-MM-DD HH24:MI:SS'),
			ap.latitude, ap.longitude, ap.location_accuracy
		FROM absensi_peserta ap
		INNER JOIN jadwal_kelas j
			ON j.id_jadwal = ap.id_jadwal
		INNER JOIN paketkelas pk
			ON pk.id_paketkelas = j.id_paketkelas
		INNER JOIN users u
			ON u.id_user = ap.id_peserta
		WHERE ap.id_jadwal = $1
		  AND ap.id_peserta = $2
		  AND ap.status = 1
	`

	var result AbsensiPesertaDTO

	err := r.DB.QueryRow(ctx, query, jadwalID, userID).Scan(
		&result.IDAbsensiPeserta, &result.IDJadwal,
		&result.IDPaketKelas, &result.NamaKelas, &result.Topik, &result.Catatan,
		&result.IDPeserta, &result.NamaPeserta, &result.NicknamePeserta,
		&result.StatusKehadiran, &result.CheckInAt,
		&result.Latitude, &result.Longitude, &result.LocationAccuracy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}


func (r *Repository) GetMentorCheckIn(
	ctx context.Context,
	jadwalID int,
) (*MentorCheckInDTO, error) {
	query := `
		SELECT
		am.id_absensi_mentor,
		am.id_jadwal,
		j.type_pertemuan,
		am.check_in_at,
		am.check_out_at,
		am.check_in_latitude,
		am.check_in_longitude,
		am.check_in_accuracy
	FROM absensi_mentor am
	INNER JOIN jadwal_kelas j
		ON j.id_jadwal = am.id_jadwal
	WHERE
		am.id_jadwal = $1
		AND am.status = 1
		AND j.status = 1
	ORDER BY am.id_absensi_mentor DESC
	LIMIT 1
	`

	var result MentorCheckInDTO

	if err := r.DB.QueryRow(ctx, query, jadwalID).Scan(
		&result.IDAbsensiMentor,
		&result.IDJadwal,
		&result.TypePertemuan,
		&result.CheckInAt,
		&result.CheckOutAt,
		&result.Latitude,
		&result.Longitude,
		&result.Accuracy,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}


func (r *Repository) GetActiveMentorCheckIn(
	ctx context.Context,
	jadwalID int,
) (*MentorCheckInDTO, error) {
	query := `
		SELECT
			id_absensi_mentor, id_jadwal,
			TO_CHAR(check_in_at, 'YYYY-MM-DD HH24:MI:SS'),
			check_in_latitude, check_in_longitude, check_in_accuracy
		FROM absensi_mentor
		WHERE id_jadwal = $1
		  AND status = 1
		  AND check_in_at IS NOT NULL
		  AND check_out_at IS NULL
		LIMIT 1
	`

	var result MentorCheckInDTO

	if err := r.DB.QueryRow(ctx, query, jadwalID).Scan(
		&result.IDAbsensiMentor,
		&result.IDJadwal,
		&result.CheckInAt,
		&result.Latitude,
		&result.Longitude,
		&result.Accuracy,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}


func (r *Repository) HasAbsensiPeserta(
	ctx context.Context,
	userID, jadwalID int,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM absensi_peserta
			WHERE id_jadwal = $1
			  AND id_peserta = $2
			  AND status = 1
		)
	`

	var exists bool
	err := r.DB.QueryRow(ctx, query, jadwalID, userID).Scan(&exists)

	return exists, err
}


func (r *Repository) InsertCheckIn(
	ctx context.Context,
	userID int,
	req CheckInPesertaRequest,
) (*AbsensiPesertaDTO, error) {

	query := `
		INSERT INTO absensi_peserta (
			id_jadwal, id_peserta, status_kehadiran,
			check_in_at, latitude, longitude, location_accuracy,
			status, created_at, updated_at
		)
		VALUES (
			$1, $2, 'HADIR',
			$3, $4, $5, $6,
			1, $3, $3
		)
		RETURNING id_absensi_peserta
	`

	now := timeutil.Now()

	var idAbsensi int

	if err := r.DB.QueryRow(
		ctx,
		query,
		req.IDJadwal,
		userID,
		now,
		req.Latitude,
		req.Longitude,
		req.LocationAccuracy,
	).Scan(&idAbsensi); err != nil {
		return nil, err
	}

	return r.GetAbsensiByID(
		ctx,
		userID,
		idAbsensi,
	)
}


func (r *Repository) GetAbsensiStatusByJadwal(
	ctx context.Context,
	userID, jadwalID int,
) (*AbsensiPesertaStatusDTO, error) {
	query := `
		SELECT
			ap.id_absensi_peserta,
			j.id_jadwal,
			ap.status_kehadiran,
			TO_CHAR(ap.check_in_at, 'YYYY-MM-DD HH24:MI:SS'),
			(ap.id_absensi_peserta IS NOT NULL)
		FROM jadwal_kelas j
		LEFT JOIN absensi_peserta ap
			ON ap.id_jadwal = j.id_jadwal
			AND ap.id_peserta = $1
			AND ap.status = 1
		WHERE j.id_jadwal = $2
		  AND j.status = 1
	`

	var result AbsensiPesertaStatusDTO

	if err := r.DB.QueryRow(ctx, query, userID, jadwalID).Scan(
		&result.IDAbsensiPeserta,
		&result.IDJadwal,
		&result.StatusKehadiran,
		&result.CheckInAt,
		&result.SudahAbsen,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}


func (r *Repository) GetAbsensiByID(
	ctx context.Context,
	userID, absensiID int,
) (*AbsensiPesertaDTO, error) {
	query := `
		SELECT
			ap.id_absensi_peserta, ap.id_jadwal,
			j.id_paketkelas, pk.nama_kelas, j.topik, j.catatan,
			ap.id_peserta, u.nama, u.nickname,
			ap.status_kehadiran,
			TO_CHAR(ap.check_in_at, 'YYYY-MM-DD HH24:MI:SS'),
			ap.latitude, ap.longitude, ap.location_accuracy
		FROM absensi_peserta ap
		INNER JOIN jadwal_kelas j
			ON j.id_jadwal = ap.id_jadwal
		INNER JOIN paketkelas pk
			ON pk.id_paketkelas = j.id_paketkelas
		INNER JOIN users u
			ON u.id_user = ap.id_peserta
		WHERE ap.id_absensi_peserta = $1
		  AND ap.id_peserta = $2
		  AND ap.status = 1
	`

	var result AbsensiPesertaDTO

	if err := r.DB.QueryRow(
		ctx, query, absensiID, userID,
	).Scan(
		&result.IDAbsensiPeserta, &result.IDJadwal,
		&result.IDPaketKelas, &result.NamaKelas, &result.Topik, &result.Catatan,
		&result.IDPeserta, &result.NamaPeserta, &result.NicknamePeserta,
		&result.StatusKehadiran, &result.CheckInAt,
		&result.Latitude, &result.Longitude, &result.LocationAccuracy,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}


func (r *Repository) GetAbsensiPeserta(
	ctx context.Context,
	userID int,
) ([]AbsensiPesertaDTO, error) {
	query := `
		SELECT
			ap.id_absensi_peserta, ap.id_jadwal,
			j.id_paketkelas, pk.nama_kelas, j.topik, j.catatan,
			ap.id_peserta, u.nama, u.nickname,
			ap.status_kehadiran,
			TO_CHAR(ap.check_in_at, 'YYYY-MM-DD HH24:MI:SS'),
			ap.latitude, ap.longitude, ap.location_accuracy
		FROM absensi_peserta ap
		INNER JOIN jadwal_kelas j
			ON j.id_jadwal = ap.id_jadwal
		INNER JOIN paketkelas pk
			ON pk.id_paketkelas = j.id_paketkelas
		INNER JOIN users u
			ON u.id_user = ap.id_peserta
		WHERE ap.id_peserta = $1
		  AND ap.status = 1
		ORDER BY ap.check_in_at DESC
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []AbsensiPesertaDTO

	for rows.Next() {
		var a AbsensiPesertaDTO

		if err := rows.Scan(
			&a.IDAbsensiPeserta, &a.IDJadwal,
			&a.IDPaketKelas, &a.NamaKelas, &a.Topik, &a.Catatan,
			&a.IDPeserta, &a.NamaPeserta, &a.NicknamePeserta,
			&a.StatusKehadiran, &a.CheckInAt,
			&a.Latitude, &a.Longitude, &a.LocationAccuracy,
		); err != nil {
			return nil, err
		}

		result = append(result, a)
	}

	return result, rows.Err()
}