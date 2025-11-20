package handlers

import (
	"github.com/Vimp17/pr-reviewer-service/internal/services"
	"github.com/gin-gonic/gin"
)

// Handlers содержит все обработчики HTTP-запросов
type Handlers struct {
	prService   *services.PRService
	teamService *services.TeamService
	userService *services.UserService
}

// NewHandlers создает новый экземпляр Handlers с указанными сервисами
func NewHandlers(
	prService *services.PRService,
	teamService *services.TeamService,
	userService *services.UserService,
) *Handlers {
	return &Handlers{
		prService:   prService,
		teamService: teamService,
		userService: userService,
	}
}

// SetupRoutes регистрирует все маршруты
func (h *Handlers) SetupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", h.healthHandler)

	// Teams endpoints
	teams := router.Group("/team")
	{
		teams.POST("/add", h.CreateTeam)
		teams.GET("/get", h.GetTeam)
	}

	// Users endpoints
	users := router.Group("/users")
	{
		users.POST("/setIsActive", h.SetUserActiveStatus)
		users.GET("/getReview", h.GetPRsForReviewer)
	}

	// PR endpoints
	pr := router.Group("/pullRequest")
	{
		pr.POST("/create", h.CreatePR)
		pr.POST("/merge", h.MergePR)
		pr.POST("/reassign", h.ReassignReviewer)
	}

	// Дополнительный эндпоинт статистики
	router.GET("/stats", h.GetStats)
}
