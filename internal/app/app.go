package app

import (
	"log"
	"net/http"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/transport"
)

type App struct {
    Addr string
}

func NewApp(adr string) *App {
    return &App{
        Addr: adr,
    }
}

func (a *App) Run() error {
    router := http.NewServeMux()

    server := &http.Server{
        Addr: a.Addr,
        Handler: transport.LogRequest(router),
    }

    log.Printf("Server started on port %s", a.Addr)
    return server.ListenAndServe()
}
