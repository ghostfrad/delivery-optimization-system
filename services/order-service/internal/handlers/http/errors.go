package http

import (
	"errors"
	"log/slog"
	"net/http"

	"order-service/internal/domain"
	"order-service/internal/service"
)

// ErrorResponse — единый формат ошибки в API.
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// writeError превращает доменную/сервисную ошибку в HTTP-ответ.
// Логирует 5xx, но клиенту отдаёт только безопасное сообщение.
func writeError(w http.ResponseWriter, logger *slog.Logger, err error) {
	status, code, message := mapError(err)

	if status >= http.StatusInternalServerError {
		logger.Error("request failed",
			"status", status,
			"code", code,
			"error", err,
		)
	}

	writeJSON(w, status, ErrorResponse{Error: message, Code: code})
}

// mapError определяет HTTP-статус, машиночитаемый код и сообщение.
func mapError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		return http.StatusNotFound, "ORDER_NOT_FOUND", "order not found"

	case errors.Is(err, service.ErrKafkaUnavailable):
		return http.StatusServiceUnavailable, "KAFKA_UNAVAILABLE", "event publishing failed, try again later"

	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"
	}
}
