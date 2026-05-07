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
		SELECT u.id_user, u.email, u.password, u.kode_pemulihan, u.role, u.status, b.status as batch_status
		FROM users u
		JOIN userbatch ub ON u.id_user = ub.id_user
		JOIN batch b ON ub.id_batch = b.id_batch
		WHERE u.email = $1
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
		&user.BatchStatus,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}