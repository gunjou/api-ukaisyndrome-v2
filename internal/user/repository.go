package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

type UserEntity struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     string
	Nickname *string
}


// ==========================
// FIND BY ID
// ==========================
func (r *Repository) FindByID(ctx context.Context, userID int) (*UserEntity, error) {

	query := `
		SELECT id_user, nama, email, password, role, nickname
		FROM users
		WHERE id_user = $1
	`

	row := r.DB.QueryRow(ctx, query, userID)

	var user UserEntity

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Nickname,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}


// ==========================
// GET USER CLASSES
// ==========================
func (r *Repository) GetUserClasses(ctx context.Context, userID int) ([]ClassDTO, error) {

	query := `
		SELECT 
			pk.id_paketkelas,
			pk.nama_kelas,
			b.id_batch,
			b.nama_batch,
			b.status
		FROM pesertakelas p
		JOIN paketkelas pk ON pk.id_paketkelas = p.id_paketkelas
		JOIN batch b ON b.id_batch = pk.id_batch
		WHERE p.id_user = $1
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []ClassDTO

	for rows.Next() {
		var c ClassDTO
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.BatchID,
			&c.BatchName,
			&c.StatusBatch,
		); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}

	return classes, nil
}


// ==========================
// GET USER MENTORSHIPS
// ==========================
func (r *Repository) GetUserMentorships(ctx context.Context, userID int) ([]MentorshipDTO, error) {

	query := `
		SELECT 
			m.id_mentorship,
			m.id_mentor,
			m.nama_mentorship,
			u.nama as mentor_name
		FROM mentorship m
		JOIN users u ON u.id_user = m.id_mentor
		WHERE m.id_peserta = $1
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mentorships []MentorshipDTO

	for rows.Next() {
		var m MentorshipDTO
		if err := rows.Scan(
			&m.ID,
			&m.MentorID,
			&m.MentorshipName,
			&m.MentorName,
		); err != nil {
			return nil, err
		}
		mentorships = append(mentorships, m)
	}

	return mentorships, nil
}


// ==========================
// UPDATE PASSWORD
// ==========================
func (r *Repository) UpdatePassword(ctx context.Context, userID int, hashed string) error {

	query := `
		UPDATE users
		SET password = $1, updated_at = NOW()
		WHERE id_user = $2
	`

	_, err := r.DB.Exec(ctx, query, hashed, userID)
	return err
}