package module

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) GetModulByUser(ctx context.Context, userID int) ([]ModulDTO, error) {

	query := `
		SELECT DISTINCT m.id_modul, m.judul
		FROM pesertakelas p
		JOIN paketkelas pk ON pk.id_paketkelas = p.id_paketkelas
		JOIN modulkelas mk ON mk.id_paketkelas = pk.id_paketkelas
		JOIN modul m ON m.id_modul = mk.id_modul
		WHERE p.id_user = $1
		AND p.status = 1
		AND pk.status = 1
		AND mk.status = 1
		AND m.status = 1
		AND m.visibility = 'open'
		ORDER BY m.judul ASC
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ModulDTO

	for rows.Next() {
		var m ModulDTO
		if err := rows.Scan(&m.ID, &m.Title); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}