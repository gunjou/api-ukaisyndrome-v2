package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	DBUrl   string
	Redis   string
	JWTSecret string 
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return Config{
		AppPort: os.Getenv("APP_PORT"),
		DBUrl:   "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
			"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") +
			"/" + os.Getenv("DB_NAME"),
		Redis: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	}
}