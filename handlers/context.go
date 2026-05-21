package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func userIDFromContext(c *gin.Context) (uint, bool) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		return 0, false
	}
	return uint(userIDFloat), true
}

func requireUserID(c *gin.Context) (uint, bool) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not identify user"})
		return 0, false
	}
	return userID, true
}
