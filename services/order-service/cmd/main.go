package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"order-service/internal/config"
	orderHTTP "order-service/internal/handlers/http"
	kafkaMock "order-service/internal/kafka/mock"
	kafkaReal "order-service/internal/kafka/real"
	repoMock "order-service/internal/repository/mock"
	repoPostgres "order-service/internal/repository/postgres"
	"order-service/internal/service"
	"order-service/pkg/database"
)

func main() {
	if err := run(); err != nil {
		slog.Error("service failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Конфиг (YAML + config.local.yaml, валидация).
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// 2. Логгер (JSON или text в зависимости от cfg.Log.Format).
	logger := newLogger(cfg.Log)
	slog.SetDefault(logger)

	logger.Info("starting order-service",
		"port", cfg.Server.Port,
		"use_mocks", cfg.UseMocks,
		"log_level", cfg.Log.Level,
		"log_format", cfg.Log.Format,
	)

	// 3. Зависимости: repo и producer.
	var (
		repo     service.OrderRepository
		producer service.MessageProducer
		cleanup  []func() error
	)

	if cfg.UseMocks {
		logger.Warn("using MOCKS for repository and Kafka")
		repo = repoMock.NewMockOrderRepository()
		producer = kafkaMock.NewMockMessageProducer()
	} else {
		// Postgres
		db, err := database.NewPostgres(cfg.Database)
		if err != nil {
			return err
		}
		cleanup = append(cleanup, db.Close)
		repo = repoPostgres.NewOrderRepository(db)
		logger.Info("connected to postgres",
			"host", cfg.Database.Host,
			"dbname", cfg.Database.DBName,
		)

		// Kafka
		kafkaProducer := kafkaReal.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.WriteTimeout)
		cleanup = append(cleanup, kafkaProducer.Close)
		producer = kafkaProducer
		logger.Info("connected to kafka",
			"brokers", cfg.Kafka.Brokers,
			"topic", cfg.Kafka.TopicOrderCreated,
		)
	}

	// 4. Очистка при выходе (в обратном порядке).
	defer func() {
		for i := len(cleanup) - 1; i >= 0; i-- {
			if err := cleanup[i](); err != nil {
				logger.Error("cleanup error", "error", err)
			}
		}
	}()

	// 5. Сервис и хендлер.
	orderService := service.NewOrderService(repo, producer, cfg.Kafka.TopicOrderCreated)
	handler := orderHTTP.NewOrderHandler(orderService, logger)

	// 6. HTTP-роутер.
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", handler.CreateOrder)
	mux.HandleFunc("/orders/", handler.GetOrder)
	mux.HandleFunc("/health", healthHandler)

	srv := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 7. Запуск сервера в горутине.
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("http server listening", "addr", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// 8. Ждём сигнал или ошибку сервера.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-quit:
		logger.Info("shutdown signal received", "signal", sig.String())
	}

	// 9. Graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server shutdown error", "error", err)
		return err
	}

	logger.Info("server stopped gracefully")
	return nil
}

// newLogger создаёт slog.Logger в зависимости от конфига.
func newLogger(cfg config.LogConfig) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLogLevel(cfg.Level)}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// healthHandler — простой health-check. Позже можно добавить проверку БД и Kafka.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
