package server

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func RespondOK(data any, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

func RespondError(err error, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	errorResponse := ErrorResponse{
		Error: err.Error(),
	}

	_ = json.NewEncoder(w).Encode(errorResponse)
}

func RespondNotFound(message string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	
	errorResponse := ErrorResponse{
		Error: message,
	}
	
	_ = json.NewEncoder(w).Encode(errorResponse)
}

func RespondBadRequest(message string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	
	errorResponse := ErrorResponse{
		Error: message,
	}
	
	_ = json.NewEncoder(w).Encode(errorResponse)
}