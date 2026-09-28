package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goose-neon-demo/models"
)

func Me(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	user, err := models.FindUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email, "created_at": user.CreatedAt})
}
