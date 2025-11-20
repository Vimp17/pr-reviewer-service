package services

import (
	"context"

	"github.com/Vimp17/pr-reviewer-service/internal/models"
	"github.com/Vimp17/pr-reviewer-service/internal/storage/postgres"
)

// UserService управляет бизнес-логикой для пользователей
type UserService struct {
	storage *postgres.Storage
}

// NewUserService создает новый сервис для работы с пользователями
func NewUserService(storage *postgres.Storage) *UserService {
	return &UserService{storage: storage}
}

// SetUserActiveStatus устанавливает флаг активности пользователя
func (s *UserService) SetUserActiveStatus(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	return s.storage.UpdateUserActiveStatus(ctx, userID, isActive)
}

// GetPRsForReviewer получает PR, где пользователь назначен ревьювером
func (s *UserService) GetPRsForReviewer(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	return s.storage.GetPRsForReviewer(ctx, userID)
}
