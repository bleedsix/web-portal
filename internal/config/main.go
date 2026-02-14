package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type Config struct {
	// Внутрішні структури для парсингу
	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`

	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"db"`

	// Приватні поля для інстансів
	logger *slog.Logger
	db     *pgxpool.Pool
}

// Load ініціалізує все одразу (Eager loading)
func Load(path string) (*Config, error) {
	cfg := &Config{}

	// 1. Налаштування Viper
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 2. Ініціалізація логера (JSON, як в logan, але стандартний)
	cfg.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// 3. Ініціалізація PGX (Connection Pool)
	// pgxpool набагато кращий за стандартний sql.DB для Postgres
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	
	// Перевірка з'єднання
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}
	cfg.db = pool

	return cfg, nil
}

// Log повертає логер
func (c *Config) Log() *slog.Logger {
	return c.logger
}

// DB повертає пул з'єднань pgx
func (c *Config) DB() *pgxpool.Pool {
	return c.db
}

// ListenerAddr формує адресу для запуску сервера
func (c *Config) ListenerAddr() string {
	return c.Server.Host + ":" + c.Server.Port
}