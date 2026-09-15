package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"order-service/internal/handlers/http/dto"
	"order-service/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Декодируем в DTO
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Валидация (вручную или через validator)
	if req.CustomerID == "" {
		http.Error(w, "customer_id is required", http.StatusBadRequest)
		return
	}
	if req.Address == "" {
		http.Error(w, "address is required", http.StatusBadRequest)
		return
	}

	// 3. DTO -> Domain (параметры) и вызов сервиса
	order, err := h.service.CreateOrder(
		r.Context(),
		req.CustomerID,
		req.Address,
		req.Lat,
		req.Lng,
		req.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Domain -> DTO
	resp := toDTOCreateOrder(order)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/orders/")
	if id == "" {
		http.Error(w, "Order ID required", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Domain -> DTO
	resp := toDTOOrder(order)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
