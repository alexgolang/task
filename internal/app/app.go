package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexgolang/ishare-task/internal/app/auth"
	"github.com/alexgolang/ishare-task/internal/app/config"
	"github.com/alexgolang/ishare-task/internal/app/db/sqlite"
	"github.com/alexgolang/ishare-task/internal/app/service"
	"github.com/alexgolang/ishare-task/internal/app/transport/httpserver"
	"github.com/alexgolang/ishare-task/internal/app/transport/httpserver/handlers"
)

type App struct {
	server *httpserver.Server
	db     *sqlite.Database
	logger *log.Logger
}

func NewApp() (*App, error) {
	cfg := config.Read()

	logger := log.New(os.Stdout, "ISHARE-TASK: ", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sqlite.NewDatabase(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	if err := db.RunMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	tokenExpiry, err := time.ParseDuration(cfg.JWTTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token expiry: %w", err)
	}

	authService, err := auth.NewJWTService(cfg.JWTPrivateKey, cfg.JWTIssuer, tokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service: %w", err)
	}

	taskService := service.NewTaskService(logger, db)

	taskHandler := handlers.NewTaskHandler(taskService)
	authHandler := handlers.NewAuthHandler(authService)

	server := httpserver.NewServer(taskHandler, authHandler, authService, cfg.Port)

	return &App{
		server: server,
		db:     db,
		logger: logger,
	}, nil
}

func (a *App) Run() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	serverErrChan := make(chan error, 1)
	go func() {
		a.logger.Printf("Starting server on port %s", a.server.GetPort())
		if err := a.server.Run(); err != nil {
			serverErrChan <- err
		}
	}()

	select {
	case sig := <-sigChan:
		a.logger.Printf("Received signal: %v. Shutting down gracefully...", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			a.logger.Printf("Server shutdown error: %v", err)
		}

		if err := a.db.Close(); err != nil {
			a.logger.Printf("Database close error: %v", err)
		}

		a.logger.Println("Shutdown complete")
		return nil

	case err := <-serverErrChan:
		a.logger.Printf("Server error: %v", err)
		return err
	}
}
