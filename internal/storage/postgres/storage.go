package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Регистрируем драйвер для database/sql
)

// Storage представляет собой хранилище данных, использующее PostgreSQL
type Storage struct {
	pool *pgxpool.Pool
}

// NewStorage создает новое хранилище с указанным DSN
func NewStorage(ctx context.Context, dsn string) (*Storage, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Настраиваем параметры пула
	poolConfig.MaxConns = 25
	poolConfig.HealthCheckPeriod = 5 * time.Minute
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{pool: pool}, nil
}

// Close закрывает соединение с базой данных
func (s *Storage) Close() {
	s.pool.Close()
}
