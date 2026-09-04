package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderHandler "order-service/internal/handlers/http"
	kafkaMock "order-service/internal/kafka/mock"
	repoMock "order-service/internal/repository/mock"
	"order-service/internal/service"
)

func main() {
	log.Println("Order Service starting...")

	// моки для тестирования
	mockRepo := repoMock.NewMockOrderRepository()
	mockProducer := kafkaMock.NewMockMessageProducer()

	orderService := service.NewOrderService(mockRepo, mockProducer)
	orderHandler := orderHandler.NewOrderHandler(orderService)

	// Настройка роутов
	http.HandleFunc("/orders", orderHandler.CreateOrder)
	http.HandleFunc("/orders/", orderHandler.GetOrder)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","mode":"mock"}`))
	})

	// Запуск сервера
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Server running on http://localhost:8080")
		log.Println("Using MOCKS for repository and Kafka")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
