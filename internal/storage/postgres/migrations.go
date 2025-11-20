package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib" // КРИТИЧЕСКИ ВАЖНО: регистрируем драйвер
	"github.com/pressly/goose/v3"
)

// ApplyMigrations применяет SQL-миграции из директории migrations
func (s *Storage) ApplyMigrations(ctx context.Context) error {
	// Получаем строку подключения из пула
	connString := s.pool.Config().ConnString()

	// Создаем *sql.DB для goose
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Определяем путь к миграциям
	migrationsDir := getMigrationsDir()

	// Применяем миграции
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// getMigrationsDir определяет путь к директории миграций
func getMigrationsDir() string {
	// Получаем текущую рабочую директорию
	cwd, err := os.Getwd()
	if err != nil {
		return "./migrations"
	}

	// Проверяем, существует ли директория миграций в текущей рабочей директории
	migrationsDir := filepath.Join(cwd, "migrations")
	if _, err := os.Stat(migrationsDir); err == nil {
		return migrationsDir
	}

	// Проверяем другие возможные расположения
	possiblePaths := []string{
		"./migrations",
		"../migrations",
		"../../migrations",
		"../../../migrations",
		"../../../../migrations",
	}

	for _, path := range possiblePaths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	// Возвращаем относительный путь по умолчанию
	return "./migrations"
}
