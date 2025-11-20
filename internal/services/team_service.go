package services

import (
	"context"
	"errors"

	"github.com/Vimp17/pr-reviewer-service/internal/models"
	"github.com/Vimp17/pr-reviewer-service/internal/storage/postgres"
)

var (
	ErrTeamExists = errors.New("TEAM_EXISTS")
)

// TeamService управляет бизнес-логикой для команд
type TeamService struct {
	storage *postgres.Storage
}

// NewTeamService создает новый сервис для работы с командами
func NewTeamService(storage *postgres.Storage) *TeamService {
	return &TeamService{storage: storage}
}

// CreateTeam создает новую команду с участниками
func (s *TeamService) CreateTeam(ctx context.Context, team models.Team) (*models.Team, error) {
	// Проверяем валидность данных
	if team.TeamName == "" {
		return nil, errors.New("TEAM_NAME_REQUIRED")
	}
	if len(team.Members) == 0 {
		return nil, errors.New("MEMBERS_REQUIRED")
	}

	// Создаем команду в хранилище
	if err := s.storage.CreateTeam(ctx, team.TeamName, team.Members); err != nil {
		if err.Error() == "TEAM_EXISTS" {
			return nil, ErrTeamExists
		}
		return nil, err
	}

	// Возвращаем созданную команду
	return s.GetTeam(ctx, team.TeamName)
}

// GetTeam получает информацию о команде
func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	return s.storage.GetTeam(ctx, teamName)
}
