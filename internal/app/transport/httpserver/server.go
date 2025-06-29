package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexgolang/ishare-task/internal/app/auth"
	"github.com/alexgolang/ishare-task/internal/app/transport/httpserver/handlers"
	"github.com/alexgolang/ishare-task/internal/app/transport/httpserver/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	taskHandler *handlers.TaskHandler
	authHandler *handlers.AuthHandler
	port        string
	srv         *http.Server
}

func NewServer(taskHandler *handlers.TaskHandler, authHandler *handlers.AuthHandler, jwtService *auth.JWTService, port string) *Server {
	router := chi.NewRouter()

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(chiMiddleware.RequestID)

	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	router.Route("/token", func(r chi.Router) {
		r.Post("/", authHandler.GetToken)
	})

	router.Route("/tasks", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		r.Post("/", taskHandler.CreateTask)
		r.Get("/", taskHandler.ListTasks)
		r.Get("/{id}", taskHandler.GetTask)
		r.Patch("/{id}", taskHandler.UpdateTask)
		r.Delete("/{id}", taskHandler.DeleteTask)
	})

	return &Server{
		taskHandler: taskHandler,
		authHandler: authHandler,
		port:        port,
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) GetPort() string {
	return s.port
}
