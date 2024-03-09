package handler

import (
	"net/http"
	"user-service/model"

	"github.com/gin-gonic/gin"
)

// Mock database
var users = map[string]model.User{
	"1": {ID: "1", Name: "John Doe", Email: "john.doe@example.com"},
}

// GetUserByID retrieves a user by ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, exists := users[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
