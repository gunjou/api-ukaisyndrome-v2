package user

import "github.com/jackc/pgx/v5/pgxpool"

type Module struct {
	Handler *Handler
}

func NewModule(db *pgxpool.Pool) *Module {

	repo := &Repository{DB: db}
	service := &Service{Repo: repo}
	handler := &Handler{Service: service}

	return &Module{
		Handler: handler,
	}
}