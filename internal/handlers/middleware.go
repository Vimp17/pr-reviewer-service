package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Обрабатываем паники
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			}})
			return
		}

		// Проверяем коды ошибок
		if c.Writer.Status() >= 400 {
			// Уже обработано в хендлерах
			return
		}
	}
}
