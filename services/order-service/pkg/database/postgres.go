package database

import (
	"database/sql"
	"fmt"

	"order-service/internal/config"

	_ "github.com/lib/pq"
)

// NewPostgres открывает соединение с PostgreSQL и настраивает пул.
// Возвращает готовый *sql.DB или ошибку, если подключиться не удалось.
func NewPostgres(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Настройка пула соединений — важно для production.
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Ping проверяет, что соединение реально работает.
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, nil
}
