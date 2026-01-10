package main

import (
	"log"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/app"
)

func main() {
    apl := app.NewApp(":9191")
    log.Println(apl.Run())
}
