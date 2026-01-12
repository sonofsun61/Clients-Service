package main

import (
	"log"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/app"
	"github.com/DATA-DOG/go-sqlmock"
)

func main() {
    db, mock, err := sqlmock.New()
    if err != nil {
        log.Fatal(err)
    }

    apl := app.NewApp(":9191", db, mock)
    log.Println(apl.Run())
}
