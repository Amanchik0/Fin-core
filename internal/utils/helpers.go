package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "user not authenticated",
		})
		return 0, false
	}

	userIDInt, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid user id format",
		})
		return 0, false
	}

	return userIDInt, true
}
