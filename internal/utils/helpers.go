package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Limit  int
	Offset int
}

var PaginationConfig = Config{
	Limit:  20,
	Offset: 0,
}

func GetPaginationParams(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit
	return limit, offset
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "user not authenticated",
		})
		return "", false
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid user id format",
		})
		return "", false
	}

	return userIDStr, true
}
