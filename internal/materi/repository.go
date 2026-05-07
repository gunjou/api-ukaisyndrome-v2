package materi

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) GetMateriByModul(ctx context.Context, userID int, modulID int, materiType *string) ([]MateriDTO, error) {

	query := `
		SELECT m.id_materi, m.id_modul, m.tipe_materi, m.judul, m.url_file, m.is_downloadable
		FROM materi m
		JOIN modul md ON md.id_modul = m.id_modul
		JOIN modulkelas mk ON mk.id_modul = md.id_modul
		JOIN paketkelas pk ON pk.id_paketkelas = mk.id_paketkelas
		JOIN pesertakelas p ON p.id_paketkelas = pk.id_paketkelas
		WHERE 
			m.id_modul = $2
			AND p.id_user = $1
			AND m.status = 1
			AND m.visibility = 'open'
	`

	args := []interface{}{userID, modulID}

	// 🔥 filter optional
	if materiType != nil {
		query += " AND m.tipe_materi = $3"
		args = append(args, *materiType)
	}

	query += " ORDER BY m.created_at ASC"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MateriDTO

	for rows.Next() {
		var m MateriDTO
		err := rows.Scan(
			&m.ID,
			&m.IDModul,
			&m.Type,
			&m.Title,
			&m.URL,
			&m.IsDownloadable,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}


func (r *Repository) GetMateriPrivateByUser(ctx context.Context, userID int, materiType *string) ([]MateriPrivateDTO, error) {

	query := `
		SELECT 
			mp.id_materi_private,
			mp.tipe_materi,
			mp.judul,
			mp.url_file,
			mp.is_downloadable
		FROM materi_private mp
		JOIN mentorship m ON m.id_mentorship = mp.id_mentorship
		WHERE 
			m.id_peserta = $1
			AND mp.status = 1
			AND mp.visibility = 'open'
	`

	args := []interface{}{userID}

	// optional filter
	if materiType != nil {
		query += " AND mp.tipe_materi = $2"
		args = append(args, *materiType)
	}

	query += " ORDER BY mp.created_at ASC"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MateriPrivateDTO

	for rows.Next() {
		var m MateriPrivateDTO
		err := rows.Scan(
			&m.ID,
			&m.Type,
			&m.Title,
			&m.URL,
			&m.IsDownloadable,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}