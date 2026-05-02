package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(dbUrl string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	return pool
}