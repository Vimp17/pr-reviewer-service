package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// Добавлен импорт services
)

type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// SetUserActiveStatus обработчик для изменения активности пользователя
func (h *Handlers) SetUserActiveStatus(c *gin.Context) {
	var req SetUserActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid user data",
		}})
		return
	}

	user, err := h.userService.SetUserActiveStatus(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "User not found",
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetPRsForReviewer обработчик для получения PR, где пользователь назначен ревьювером
func (h *Handlers) GetPRsForReviewer(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "user_id is required",
		}})
		return
	}

	prs, err := h.userService.GetPRsForReviewer(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
