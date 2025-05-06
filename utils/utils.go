package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse envoie une réponse JSON standard pour les succès
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// ErrorResponse envoie une réponse JSON standard pour les erreurs
func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error": gin.H{
			"message": message,
		},
	})
}
