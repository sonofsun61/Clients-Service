package app

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/database"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/service"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/transport"
)

const (
	shutdownTimeout = 10
)

type App struct {
	Addr string
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func NewApp(adr string, db *sql.DB, mock sqlmock.Sqlmock) *App {
	return &App{
		Addr: adr,
		db:   db,
		mock: mock,
	}
}

func (a *App) Run() {
	router := http.NewServeMux()

	mocker := database.NewMocker(a.mock)
	profileRepository := database.NewProfileRepository(a.db, mocker)
	authRepository := database.NewAuthRepository(a.db, mocker)
	profileService := service.NewProfileService(profileRepository)
	authService := service.NewAuthService(authRepository)
	handler := transport.NewHandler(profileService, authService)

	handler.RegisterRoutes(router)

	server := &http.Server{
		Addr:    a.Addr,
		Handler: transport.LogRequest(router),
	}

	go func() {
		log.Printf("Server started on port %s", a.Addr)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("Error during server starting:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Error during server shutting down %v.\n", err)
	}

	log.Println("Server was shut down")
}
