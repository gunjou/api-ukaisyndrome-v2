package main

import (
	"log"
	"os"

	_ "api-ukaisyndrome-v2/docs"

	"api-ukaisyndrome-v2/docs"
	"api-ukaisyndrome-v2/internal/app"
	"api-ukaisyndrome-v2/pkg/config"
)

// @title UKAI Syndrome API V2
// @version 2.0
// @description API Documentation
// @BasePath /api/v2

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	cfg := config.LoadConfig()

	// ==========================
	// SWAGGER CONFIG (DYNAMIC)
	// ==========================
	swaggerHost := os.Getenv("SWAGGER_HOST")
	swaggerScheme := os.Getenv("SWAGGER_SCHEME")

	if swaggerHost != "" {
		docs.SwaggerInfo.Host = swaggerHost
	}

	if swaggerScheme != "" {
		docs.SwaggerInfo.Schemes = []string{swaggerScheme}
	}

	// ==========================
	// INIT APP
	// ==========================
	application := app.NewApp(cfg)

	log.Println("🚀 Server running on port:", cfg.AppPort)
	log.Fatal(application.Fiber.Listen(":" + cfg.AppPort))
}