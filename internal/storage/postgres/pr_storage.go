package postgres

import (
	"context"
	"errors"

	"github.com/Vimp17/pr-reviewer-service/internal/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = errors.New("not found")
)

// CheckPRExists проверяет существование PR
func (s *Storage) CheckPRExists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)
	`, prID).Scan(&exists)
	return exists, err
}

// CreatePR создает новый PR
func (s *Storage) CreatePR(ctx context.Context, pr models.PullRequest) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO pull_requests (
			pull_request_id, pull_request_name, author_id, status, 
			reviewer1_id, reviewer2_id
		) VALUES ($1, $2, $3, $4, $5, $6)
	`,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
		getReviewer(pr.AssignedReviewers, 0),
		getReviewer(pr.AssignedReviewers, 1),
	)
	return err
}

// GetPR получает информацию о PR
func (s *Storage) GetPR(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr := &models.PullRequest{}
	var reviewer1, reviewer2 *string

	err := s.pool.QueryRow(ctx, `
		SELECT 
			pull_request_id, pull_request_name, author_id, status,
			reviewer1_id, reviewer2_id, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&reviewer1,
		&reviewer2,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Преобразуем указатели в слайс
	pr.AssignedReviewers = make([]string, 0, 2)
	if reviewer1 != nil {
		pr.AssignedReviewers = append(pr.AssignedReviewers, *reviewer1)
	}
	if reviewer2 != nil {
		pr.AssignedReviewers = append(pr.AssignedReviewers, *reviewer2)
	}

	return pr, nil
}

// UpdatePRReviewers обновляет список ревьюеров для PR
func (s *Storage) UpdatePRReviewers(ctx context.Context, prID string, reviewers []string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE pull_requests
		SET reviewer1_id = $2, reviewer2_id = $3
		WHERE pull_request_id = $1
	`,
		prID,
		getReviewer(reviewers, 0),
		getReviewer(reviewers, 1),
	)
	return err
}

// MergePR помечает PR как MERGED
// MergePR помечает PR как MERGED
func (s *Storage) MergePR(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr := &models.PullRequest{}
	var reviewer1, reviewer2 *string

	err := s.pool.QueryRow(ctx, `
        UPDATE pull_requests 
        SET status = 'MERGED', merged_at = NOW()
        WHERE pull_request_id = $1
        RETURNING 
            pull_request_id, pull_request_name, author_id, status,
            reviewer1_id, reviewer2_id, created_at, merged_at
    `, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&reviewer1,
		&reviewer2,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Преобразуем указатели в слайс
	pr.AssignedReviewers = make([]string, 0, 2)
	if reviewer1 != nil {
		pr.AssignedReviewers = append(pr.AssignedReviewers, *reviewer1)
	}
	if reviewer2 != nil {
		pr.AssignedReviewers = append(pr.AssignedReviewers, *reviewer2)
	}

	return pr, nil
}

// GetAssignmentStats возвращает статистику по назначениям
func (s *Storage) GetAssignmentStats(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT reviewer_id, COUNT(*)
		FROM (
			SELECT reviewer1_id AS reviewer_id FROM pull_requests WHERE reviewer1_id IS NOT NULL
			UNION ALL
			SELECT reviewer2_id FROM pull_requests WHERE reviewer2_id IS NOT NULL
		) AS reviewers
		GROUP BY reviewer_id
	`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var userID string
		var count int
		if err := rows.Scan(&userID, &count); err != nil {
			return nil, err
		}
		stats[userID] = count
	}

	return stats, rows.Err()
}

// Вспомогательные функции

func getReviewer(reviewers []string, index int) *string {
	if index < len(reviewers) {
		return &reviewers[index]
	}
	return nil
}
