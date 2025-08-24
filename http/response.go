package http

import (
	"encoding/json"
	"net/http"

	"github.com/rsdlab-dk/tft-core/logger"
	"go.uber.org/zap"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, data interface{}, log *logger.Logger, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Success: true,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.WithContext(r.Context()).Error("failed to encode json response",
			zap.Error(err))
		WriteError(w, "ENCODING_ERROR", "Failed to encode response", http.StatusInternalServerError, log, r)
		return
	}
}

func WriteError(w http.ResponseWriter, code, message string, statusCode int, log *logger.Logger, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.WithContext(r.Context()).Error("failed to encode error response",
			zap.Error(err))
	}
}

func WriteNotFound(w http.ResponseWriter, message string, log *logger.Logger, r *http.Request) {
	WriteError(w, "NOT_FOUND", message, http.StatusNotFound, log, r)
}

func WriteBadRequest(w http.ResponseWriter, message string, log *logger.Logger, r *http.Request) {
	WriteError(w, "BAD_REQUEST", message, http.StatusBadRequest, log, r)
}

func WriteInternalError(w http.ResponseWriter, log *logger.Logger, r *http.Request) {
	WriteError(w, "INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError, log, r)
}
