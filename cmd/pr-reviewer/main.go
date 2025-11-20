package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vimp17/pr-reviewer-service/internal/handlers"
	"github.com/Vimp17/pr-reviewer-service/internal/services"
	"github.com/Vimp17/pr-reviewer-service/internal/storage/postgres"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	// 1. Подключаемся к БД
	dsn := os.Getenv("DB_CONN_STRING")
	if dsn == "" {
		log.Fatal("DB_CONN_STRING environment variable is required")
	}

	storage, err := postgres.NewStorage(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer storage.Close()

	// 2. Применяем миграции
	if err := storage.ApplyMigrations(ctx); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// 3. Инициализируем сервисы
	prService := services.NewPRService(storage)
	teamService := services.NewTeamService(storage)
	userService := services.NewUserService(storage)

	// 4. Настраиваем роутер
	router := gin.Default()

	// Создаем обработчики
	h := handlers.NewHandlers(prService, teamService, userService)

	// Регистрируем маршруты
	h.SetupRoutes(router)

	// 5. Запускаем сервер
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}
