package services

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/Vimp17/pr-reviewer-service/internal/models"
	"github.com/Vimp17/pr-reviewer-service/internal/storage/postgres"
)

var (
	ErrPRExists       = errors.New("PR_EXISTS")
	ErrAuthorNotFound = errors.New("AUTHOR_NOT_FOUND")
	ErrPRMerged       = errors.New("PR_MERGED")
	ErrNotAssigned    = errors.New("NOT_ASSIGNED")
	ErrNoCandidate    = errors.New("NO_CANDIDATE")
	ErrNotFound       = errors.New("NOT_FOUND")
)

// PRService управляет бизнес-логикой для Pull Requests
type PRService struct {
	storage *postgres.Storage
}

// NewPRService создает новый сервис для работы с PR
func NewPRService(storage *postgres.Storage) *PRService {
	return &PRService{storage: storage}
}

// CreatePR создает новый PR и назначает ревьюеров
func (s *PRService) CreatePR(ctx context.Context, pr models.PullRequest) (*models.PullRequest, error) {
	// Проверяем существование PR
	exists, err := s.storage.CheckPRExists(ctx, pr.PullRequestID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPRExists
	}

	// Получаем автора
	author, err := s.storage.GetUser(ctx, pr.AuthorID)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, ErrAuthorNotFound
		}
		return nil, err
	}

	// Получаем активных членов команды (кроме автора)
	candidates, err := s.storage.GetActiveTeamMembers(ctx, author.TeamName, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	// Назначаем до 2 ревьюеров
	reviewers := selectReviewers(candidates)
	pr.AssignedReviewers = reviewers
	pr.Status = "OPEN"

	// Сохраняем в БД
	if err := s.storage.CreatePR(ctx, pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

// MergePR помечает PR как MERGED
func (s *PRService) MergePR(ctx context.Context, prID string) (*models.PullRequest, error) {
	// Получаем текущий статус PR
	pr, err := s.storage.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}

	// Если уже merged, возвращаем текущее состояние
	if pr.Status == "MERGED" {
		return pr, nil
	}

	// Обновляем статус
	mergedPR, err := s.storage.MergePR(ctx, prID)
	if err != nil {
		return nil, err
	}

	return mergedPR, nil
}

// ReassignReviewer заменяет одного ревьюера на другого из его команды
func (s *PRService) ReassignReviewer(
	ctx context.Context,
	prID, oldUserID string,
) (*models.PullRequest, string, error) {
	// Получаем PR
	pr, err := s.storage.GetPR(ctx, prID)
	if err != nil {
		return nil, "", err
	}

	// Проверяем статус
	if pr.Status == "MERGED" {
		return nil, "", ErrPRMerged
	}

	// Проверяем, назначен ли пользователь
	if !contains(pr.AssignedReviewers, oldUserID) {
		return nil, "", ErrNotAssigned
	}

	// Получаем информацию о заменяемом ревьюере
	reviewer, err := s.storage.GetUser(ctx, oldUserID)
	if err != nil {
		return nil, "", err
	}

	// Ищем замену в его команде
	candidates, err := s.storage.GetActiveTeamMembers(ctx, reviewer.TeamName, oldUserID)
	if err != nil {
		return nil, "", err
	}

	if len(candidates) == 0 {
		return nil, "", ErrNoCandidate
	}

	// Случайный выбор
	rand.Seed(time.Now().UnixNano())
	newReviewer := candidates[rand.Intn(len(candidates))]

	// Обновляем назначения
	pr.AssignedReviewers = replaceReviewer(pr.AssignedReviewers, oldUserID, newReviewer)
	if err := s.storage.UpdatePRReviewers(ctx, prID, pr.AssignedReviewers); err != nil {
		return nil, "", err
	}

	return pr, newReviewer, nil
}

// GetAssignmentStats возвращает статистику по назначениям
func (s *PRService) GetAssignmentStats(ctx context.Context) (map[string]int, error) {
	return s.storage.GetAssignmentStats(ctx)
}

// Вспомогательные функции

func selectReviewers(candidates []string) []string {
	if len(candidates) == 0 {
		return []string{}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if len(candidates) >= 2 {
		return candidates[:2]
	}
	return candidates
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func replaceReviewer(reviewers []string, oldID, newID string) []string {
	for i, id := range reviewers {
		if id == oldID {
			reviewers[i] = newID
			break
		}
	}
	return reviewers
}
