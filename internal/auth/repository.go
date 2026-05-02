package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*UserAuthEntity, error) {
	query := `
		SELECT id_user, email, password, kode_pemulihan, role, status
		FROM users
		WHERE email = $1
	`

	row := r.DB.QueryRow(ctx, query, email)

	var user UserAuthEntity
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.KodePemulihan,
		&user.Role,
		&user.Status,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}