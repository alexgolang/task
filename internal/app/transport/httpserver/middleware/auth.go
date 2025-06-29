package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexgolang/ishare-task/internal/app/auth"
	"github.com/alexgolang/ishare-task/internal/app/common/server"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	jwtService *auth.JWTService
}

func NewAuthMiddleware(jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			server.RespondBadRequest("Authorization header required", w, r)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			server.RespondBadRequest("Invalid authorization header format. Expected 'Bearer <token>'", w, r)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			server.RespondBadRequest("Token required", w, r)
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			server.RespondBadRequest("Invalid token: "+err.Error(), w, r)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
