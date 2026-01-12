package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/database"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/service"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/transport"
	"github.com/DATA-DOG/go-sqlmock"
)

type App struct {
    Addr string
    db *sql.DB
    mock sqlmock.Sqlmock
}

func NewApp(adr string, db *sql.DB, mock sqlmock.Sqlmock) *App {
    return &App{
        Addr: adr,
        db: db,
        mock: mock,
    }
}

func (a *App) Run() error {
    router := http.NewServeMux()

    mocker := database.NewMocker(a.mock)
    repository := database.NewRepository(a.db, mocker)
    service := service.NewService(repository)
    handler := transport.NewHandler(service)

    handler.RegisterRoutes(router)

    server := &http.Server{
        Addr: a.Addr,
        Handler: transport.LogRequest(router),
    }

    log.Printf("Server started on port %s", a.Addr)
    return server.ListenAndServe()
}
