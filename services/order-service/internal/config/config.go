package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config — корневая структура конфигурации.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Kafka    KafkaConfig    `yaml:"kafka"`
	Log      LogConfig      `yaml:"log"`
	UseMocks bool           `yaml:"use_mocks"`
}

// ServerConfig — настройки HTTP-сервера.
type ServerConfig struct {
	Port            string        `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// DatabaseConfig — настройки подключения к PostgreSQL.
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// DSN собирает строку подключения для lib/pq.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// KafkaConfig — настройки Kafka-продюсера.
type KafkaConfig struct {
	Brokers           []string      `yaml:"brokers"`
	TopicOrderCreated string        `yaml:"topic_order_created"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	RequiredAcks      string        `yaml:"required_acks"`
}

// LogConfig — настройки логирования.
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	basePath := getEnvOrDefault("CONFIG_PATH", "configs/config.yaml")
	if err := loadYAML(basePath, cfg, true); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	return cfg, nil
}

// loadYAML читает YAML и мержит в cfg.
// required=true → файл обязателен; required=false → файл опционален.
func loadYAML(path string, cfg *Config, required bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && !required {
			return nil
		}
		return fmt.Errorf("read config %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("unmarshal config %s: %w", path, err)
	}
	return nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Validate проверяет обязательные поля и допустимые значения.
// Вызывается после загрузки всех конфигов.
func (c *Config) Validate() error {
	// --- Server ---
	if c.Server.Port == "" {
		return errors.New("server.port is required")
	}
	if c.Server.ReadTimeout <= 0 {
		return errors.New("server.read_timeout must be > 0")
	}
	if c.Server.WriteTimeout <= 0 {
		return errors.New("server.write_timeout must be > 0")
	}
	if c.Server.ShutdownTimeout <= 0 {
		return errors.New("server.shutdown_timeout must be > 0")
	}

	// --- Database и Kafka проверяем, только если не работаем на моках ---
	if !c.UseMocks {
		if c.Database.Host == "" {
			return errors.New("database.host is required")
		}
		if c.Database.Port == "" {
			return errors.New("database.port is required")
		}
		if c.Database.User == "" {
			return errors.New("database.user is required")
		}
		if c.Database.Password == "" {
			return errors.New("database.password is required (set in config.local.yaml)")
		}
		if c.Database.DBName == "" {
			return errors.New("database.dbname is required")
		}
		if c.Database.MaxOpenConns <= 0 {
			return errors.New("database.max_open_conns must be > 0")
		}
		if c.Database.MaxIdleConns < 0 {
			return errors.New("database.max_idle_conns must be >= 0")
		}
		if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
			return errors.New("database.max_idle_conns cannot exceed max_open_conns")
		}
		if c.Database.ConnMaxLifetime <= 0 {
			return errors.New("database.conn_max_lifetime must be > 0")
		}

		if len(c.Kafka.Brokers) == 0 {
			return errors.New("kafka.brokers is required")
		}
		if c.Kafka.TopicOrderCreated == "" {
			return errors.New("kafka.topic_order_created is required")
		}
		if c.Kafka.WriteTimeout <= 0 {
			return errors.New("kafka.write_timeout must be > 0")
		}
	}

	// --- Log ---
	switch c.Log.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("log.level must be one of: debug|info|warn|error, got %q", c.Log.Level)
	}
	switch c.Log.Format {
	case "json", "text":
	default:
		return fmt.Errorf("log.format must be one of: json|text, got %q", c.Log.Format)
	}

	return nil
}
