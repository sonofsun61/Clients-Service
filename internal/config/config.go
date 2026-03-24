package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DB        *sql.DB
	Mocker    sqlmock.Sqlmock
	JWTSecret string
}

func Load() *Config {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Println("No .env file found in configs/, using system environment variables")
	}
	db, mocker, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		Port:      getEnv("PORT", "9191"),
		DB:        db,
		Mocker:    mocker,
		JWTSecret: getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
