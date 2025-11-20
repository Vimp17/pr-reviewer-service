package handlers

import (
	"net/http"

	"github.com/Vimp17/pr-reviewer-service/internal/models"
	"github.com/Vimp17/pr-reviewer-service/internal/services" // Добавлен импорт services
	"github.com/gin-gonic/gin"
)

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_user_id" binding:"required"`
}

// CreatePR обработчик для создания PR
func (h *Handlers) CreatePR(c *gin.Context) {
	var req CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request body",
		}})
		return
	}

	// Преобразуем в модель
	pr := models.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
		Status:          "OPEN",
	}

	// Вызываем сервис
	createdPR, err := h.prService.CreatePR(c.Request.Context(), pr)
	if err != nil {
		switch {
		case err == services.ErrPRExists:
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{
				"code":    "PR_EXISTS",
				"message": "PR id already exists",
			}})
		case err == services.ErrAuthorNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Author not found",
			}})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pr": createdPR})
}

// MergePR обработчик для пометки PR как MERGED
func (h *Handlers) MergePR(c *gin.Context) {
	var req MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid PR ID",
		}})
		return
	}

	pr, err := h.prService.MergePR(c.Request.Context(), req.PullRequestID)
	if err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "PR not found",
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

// ReassignReviewer обработчик для переназначения ревьюера
func (h *Handlers) ReassignReviewer(c *gin.Context) {
	var req ReassignReviewerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request",
		}})
		return
	}

	pr, newReviewer, err := h.prService.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		switch {
		case err == services.ErrPRMerged:
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{
				"code":    "PR_MERGED",
				"message": "cannot reassign on merged PR",
			}})
		case err == services.ErrNotAssigned:
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{
				"code":    "NOT_ASSIGNED",
				"message": "reviewer is not assigned to this PR",
			}})
		case err == services.ErrNoCandidate:
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{
				"code":    "NO_CANDIDATE",
				"message": "no active replacement candidate in team",
			}})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr, "replaced_by": newReviewer})
}
