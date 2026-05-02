package app

import (
	"api-ukaisyndrome-v2/pkg/config"
	"api-ukaisyndrome-v2/pkg/database"
	"api-ukaisyndrome-v2/pkg/middleware"
	redisPkg "api-ukaisyndrome-v2/pkg/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type App struct {
	Fiber *fiber.App
}

func NewApp(cfg config.Config) *App {

	// ==========================
	// INIT FIBER
	// ==========================
	app := fiber.New(fiber.Config{
		AppName: "UKAI Syndrome API V2",
	})

	// ==========================
	// GLOBAL MIDDLEWARE (URUTAN PENTING)
	// ==========================

	// 1. CORS (PALING ATAS)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // dev mode
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// 2. Logger (biar kelihatan request masuk)
	app.Use(logger.New())

	// 3. Request ID (trace request)
	app.Use(middleware.RequestID())

	// 4. Response Time
	app.Use(middleware.ResponseTime())

	// ==========================
	// INIT INFRASTRUCTURE
	// ==========================
	db := database.NewPostgres(cfg.DBUrl)
	rdb := redisPkg.NewRedis(cfg.Redis)

	// ==========================
	// SWAGGER
	// ==========================
	app.Get("/docs/*", fiberSwagger.WrapHandler)

	// ==========================
	// API GROUP
	// ==========================
	api := app.Group("/api/v2")

	// ==========================
	// REGISTER MODULES
	// ==========================
	registerModules(api, db, rdb, cfg)

	return &App{
		Fiber: app,
	}
}