package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
}

func Load() *Config {
	err := godotenv.Load("../../configs/.env")
	if err != nil {
		log.Println(err)
	}
	return &Config{
		Port:      getEnv("PORT", "9191"),
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
