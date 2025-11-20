package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// Добавлен импорт services
)

// GetStats обработчик для получения статистики назначений
func (h *Handlers) GetStats(c *gin.Context) {
	stats, err := h.prService.GetAssignmentStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// healthHandler обработчик для проверки работоспособности
func (h *Handlers) healthHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}
