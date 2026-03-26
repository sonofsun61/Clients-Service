package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/config"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/repository"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/service"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/transport"
)

const (
	shutdownTimeout = 10
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mocker := repository.NewMocker(a.cfg.Mocker)
	profileRepo := repository.NewProfileRepository(a.cfg.DB, mocker)
	authRepo := repository.NewAuthRepository(a.cfg.DB, mocker)

	profileService := service.NewProfileService(profileRepo) // можно добавить логгер
	authService := service.NewAuthService(authRepo, logger)

	mainRouter := http.NewServeMux()

	h := transport.NewHandler(profileService, authService, logger)

	h.RegisterPublicRoutes(mainRouter)

	protectedMux := http.NewServeMux()
	h.RegisterProtectedRoutes(protectedMux)

	mainRouter.Handle("/", transport.AuthMiddleware(protectedMux))

	finalHandler := transport.LogRequest(logger, mainRouter)

	server := &http.Server{
		Addr:     ":" + a.cfg.Port,
		Handler:  finalHandler,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	go func() {
		logger.Info("Server started", slog.String("addr", a.cfg.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server startup failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed", "error", err)
	}

	logger.Info("Server stopped")
}
