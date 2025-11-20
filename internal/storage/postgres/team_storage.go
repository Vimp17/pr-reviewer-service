package postgres

import (
	"context"
	"errors"

	"github.com/Vimp17/pr-reviewer-service/internal/models"
	"github.com/jackc/pgx/v5"
)

// CheckTeamExists проверяет существование команды
func (s *Storage) CheckTeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)
	`, teamName).Scan(&exists)
	return exists, err
}

// CreateTeam создает новую команду с участниками
func (s *Storage) CreateTeam(ctx context.Context, teamName string, members []models.User) error {
	// Проверяем существование команды
	exists, err := s.CheckTeamExists(ctx, teamName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("TEAM_EXISTS")
	}

	// Начинаем транзакцию
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Создаем команду
	if _, err := tx.Exec(ctx, "INSERT INTO teams (team_name) VALUES ($1)", teamName); err != nil {
		return err
	}

	// Обрабатываем участников
	for _, member := range members {
		// Обновляем или создаем пользователя
		_, err := tx.Exec(ctx, `
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) 
			DO UPDATE SET team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active
		`, member.UserID, member.Username, teamName, member.IsActive)

		if err != nil {
			return err
		}
	}

	// Фиксируем изменения
	return tx.Commit(ctx)
}

// GetTeam получает информацию о команде
func (s *Storage) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	// Проверяем существование команды
	exists, err := s.CheckTeamExists(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("NOT_FOUND")
	}

	// Получаем участников
	rows, err := s.pool.Query(ctx, `
		SELECT user_id, username, is_active 
		FROM users 
		WHERE team_name = $1
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.IsActive); err != nil {
			return nil, err
		}
		user.TeamName = teamName // Добавляем team_name в модель
		members = append(members, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}
