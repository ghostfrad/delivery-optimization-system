package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"order-service/internal/handlers/http/dto"
	"order-service/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
	logger  *slog.Logger
}

func NewOrderHandler(service *service.OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{
		service: service,
		logger:  logger,
	}
}

// CreateOrder — POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
			Code:  "METHOD_NOT_ALLOWED",
		})
		return
	}

	// 1. Декодируем DTO
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
			Code:  "INVALID_BODY",
		})
		return
	}

	// 2. Валидация (пока вручную)
	if req.CustomerID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "customer_id is required",
			Code:  "VALIDATION_ERROR",
		})
		return
	}
	if req.Address == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "address is required",
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	// 3. Вызов сервиса
	order, err := h.service.CreateOrder(
		r.Context(),
		req.CustomerID,
		req.Address,
		req.Lat,
		req.Lng,
		req.Description,
	)
	if err != nil {
		writeError(w, h.logger, err)
		return
	}

	// 4. Domain -> DTO
	resp := toDTOCreateOrder(order)
	writeJSON(w, http.StatusCreated, resp)
}

// GetOrder — GET /orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
			Code:  "METHOD_NOT_ALLOWED",
		})
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/orders/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "order id is required",
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		writeError(w, h.logger, err)
		return
	}

	resp := toDTOOrder(order)
	writeJSON(w, http.StatusOK, resp)
}
