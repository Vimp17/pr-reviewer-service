package postgres

import (
	"context"
	"errors"

	"github.com/Vimp17/pr-reviewer-service/internal/models"
	"github.com/jackc/pgx/v5"
)

// UpdateUserActiveStatus обновляет статус активности пользователя
func (s *Storage) UpdateUserActiveStatus(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	_, err := s.pool.Exec(ctx, `
		UPDATE users 
		SET is_active = $1 
		WHERE user_id = $2
	`, isActive, userID)
	if err != nil {
		return nil, err
	}

	// Получаем обновленного пользователя
	var user models.User
	err = s.pool.QueryRow(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`, userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUser получает пользователя по ID
func (s *Storage) GetUser(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := s.pool.QueryRow(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`, userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetActiveTeamMembers возвращает активных членов команды, исключая указанного пользователя
func (s *Storage) GetActiveTeamMembers(ctx context.Context, teamName, excludeUserID string) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT user_id
		FROM users
		WHERE team_name = $1 AND is_active = true AND user_id != $2
	`, teamName, excludeUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}

	return users, rows.Err()
}

// GetPRsForReviewer возвращает PR, где пользователь назначен ревьювером
func (s *Storage) GetPRsForReviewer(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT 
			pull_request_id, 
			pull_request_name, 
			author_id, 
			status
		FROM pull_requests
		WHERE reviewer1_id = $1 OR reviewer2_id = $1
		AND status = 'OPEN'
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
		); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, rows.Err()
}
