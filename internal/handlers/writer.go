package handlers

import (
	"app/internal/repository"
	"app/internal/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var ErrInternal = errors.New("internal error")

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		slog.Error("RespondJSON: data is nil")
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response",
			slog.Any("error", err),
		)
		RespondError(w, ErrInternal)
		return
	}

}

func RespondError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	statusCode, errorCode, message := MapErrorToHTTP(err)
	errResp := ErrorResponse{
		Error: ErrorDetail{
			Code:    errorCode,
			Message: message,
		},
	}
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		slog.Error("Failed to encode error response",
			slog.Any("error", err),
		)
		return
	}

}

func MapErrorToHTTP(err error) (statusCode int, errorCode string, message string) {
	if err == nil {
		return http.StatusOK, "", ""
	}

	switch {

	case errors.Is(err, ErrInvalidRequest):

		return http.StatusBadRequest, "INVALID_REQUEST", "invalid request"
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound, "NOT_FOUND", "resource not found"
	case errors.Is(err, service.ErrInvalidDates):
		return http.StatusBadRequest, "INVALID_REQUEST", "Invalid dates"
	case errors.Is(err, service.ErrInvalidPageArgs):
		return http.StatusBadRequest, "INVALID_REQUEST", "invalid page args"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"
	}

}
