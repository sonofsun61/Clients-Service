package main

import (
	"log"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/app"
)

const (
    defaultPort = ":9191"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    port, exists := os.LookupEnv("PORT")
    if !exists {
        port = defaultPort
    }

    db, mock, err := sqlmock.New()
    if err != nil {
        log.Fatal(err)
    }

    apl := app.NewApp(port, db, mock)
    apl.Run()
}
