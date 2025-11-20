package handlers

import (
	"net/http"

	"github.com/Vimp17/pr-reviewer-service/internal/models"
	"github.com/Vimp17/pr-reviewer-service/internal/services" // Добавлен импорт services
	"github.com/gin-gonic/gin"
)

type CreateTeamRequest struct {
	TeamName string          `json:"team_name" binding:"required"`
	Members  []TeamMemberDTO `json:"members" binding:"required,min=1"`
}

type TeamMemberDTO struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// CreateTeam обработчик для создания команды
func (h *Handlers) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid team data",
		}})
		return
	}

	// Преобразуем DTO в модель
	members := make([]models.User, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, models.User{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	team := models.Team{
		TeamName: req.TeamName,
		Members:  members,
	}

	// Вызываем сервис
	createdTeam, err := h.teamService.CreateTeam(c.Request.Context(), team)
	if err != nil {
		if err == services.ErrTeamExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
				"code":    "TEAM_EXISTS",
				"message": "team_name already exists",
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": createdTeam})
}

// GetTeam обработчик для получения информации о команде
func (h *Handlers) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "team_name is required",
		}})
		return
	}

	team, err := h.teamService.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Team not found",
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, team)
}
