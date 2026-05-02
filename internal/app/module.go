package app

import (
	"api-ukaisyndrome-v2/internal/auth"
	"api-ukaisyndrome-v2/internal/materi"
	"api-ukaisyndrome-v2/internal/module"
	"api-ukaisyndrome-v2/internal/user"
	"api-ukaisyndrome-v2/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func registerModules(r fiber.Router, db *pgxpool.Pool, rdb *redis.Client, cfg config.Config) {

	// AUTH
	authRepo := &auth.Repository{DB: db}
	authService := &auth.Service{
		Repo:      authRepo,
		Redis:     rdb,
		JWTSecret: cfg.JWTSecret,
	}
	authHandler := &auth.Handler{Service: authService}
	auth.RegisterRoutes(r, authHandler)

	// PROTECTED GROUP
	protected := r.Group("", auth.AuthMiddleware(rdb, cfg.JWTSecret))

	// USER
	userModule := user.NewModule(db)
	user.RegisterRoutes(protected, userModule.Handler)

	// MODUL
	moduleRepo := &module.Repository{DB: db}
	moduleService := &module.Service{Repo: moduleRepo}
	moduleHandler := &module.Handler{Service: moduleService}

	module.RegisterRoutes(protected, moduleHandler)

	// MATERI
	materiRepo := &materi.Repository{DB: db}
	materiService := &materi.Service{Repo: materiRepo}
	materiHandler := &materi.Handler{Service: materiService}

	materi.RegisterRoutes(protected, materiHandler)
}