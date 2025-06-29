package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexgolang/ishare-task/internal/app/auth"
	"github.com/alexgolang/ishare-task/internal/app/common/server"
	"github.com/alexgolang/ishare-task/internal/app/domain"
)

type AuthHandler struct {
	JWTService *auth.JWTService
}

func NewAuthHandler(jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		JWTService: jwtService,
	}
}

func (h *AuthHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		server.RespondError(err, w, r)
		return
	}

	grantType := r.FormValue("grant_type")
	if grantType != "client_credentials" {
		server.RespondError(errors.New("unsupported grant_type"), w, r)
		return
	}

	clientAssertion := r.FormValue("client_assertion")
	if clientAssertion == "" {
		server.RespondError(errors.New("client_assertion is required"), w, r)
		return
	}

	clientAssertionType := r.FormValue("client_assertion_type")
	if clientAssertionType == "" {
		server.RespondError(errors.New("client_assertion_type is required"), w, r)
	}

	claims, err := h.JWTService.ValidateClientAssertion(clientAssertion, clientAssertionType)
	if err != nil {
		server.RespondBadRequest("Invalid client assertion: "+err.Error(), w, r)
		return
	}

	clientID, ok := claims["sub"].(string)
	if !ok || clientID == "" {
		server.RespondBadRequest("Invalid or missing 'sub' claim in client assertion", w, r)
		return
	}

	accessToken, err := h.JWTService.CreateAccessToken(clientID)
	if err != nil {
		server.RespondError(fmt.Errorf("failed to create access token: %w", err), w, r)
		return
	}

	resp := domain.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	server.RespondOK(resp, w, r)
}
