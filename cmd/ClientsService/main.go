package main

import (
	"github.com/AI-Hackathon-2026/Clients-Service/internal/app"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/config"
)

func main() {
	cfg := config.Load()
	apl := app.NewApp(cfg)
	apl.Run()
}
