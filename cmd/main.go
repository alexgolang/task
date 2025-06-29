package main

import (
	"fmt"
	"log"

	"github.com/alexgolang/ishare-task/internal/app"
	_ "github.com/alexgolang/ishare-task/docs"
)

// @title Task Management API
// @version 1.0
// @description A REST API for task management with CRUD operations
// @host localhost:8080
// @BasePath /

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}

func run() error {
	app, err := app.NewApp()
	if err != nil {
		return fmt.Errorf("failed to create app: %v", err)
	}

	return app.Run()
}
