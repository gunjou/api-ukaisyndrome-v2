package materi

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) GetMateriByModul(ctx context.Context, modulID int, materiType *string) ([]MateriDTO, error) {

	query := `
		SELECT id_materi, id_modul, tipe_materi, judul, url_file, is_downloadable
		FROM materi
		WHERE id_modul = $1
		AND status = 1
		AND visibility = 'open'
	`

	args := []interface{}{modulID}

	// 🔥 filter optional
	if materiType != nil {
		query += " AND tipe_materi = $2"
		args = append(args, *materiType)
	}

	query += " ORDER BY created_at ASC"

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