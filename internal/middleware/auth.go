package middleware

import (
	"justTest/internal/infrastructure/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authClient *auth.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"ereor":   "Authorization header is required",
			})
			c.Abort()
			return
		}

		user, err := authClient.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "incalid to expired token",
				"details": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Set("token", token)
		c.Next()

	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func OptionalAuthMiddleware(authClient *auth.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Пытаемся получить токен
		token, err := c.Cookie("token")
		if err != nil {
			// Нет токена - продолжаем без аутентификации
			c.Next()
			return
		}

		// 3. Валидируем токен
		user, err := authClient.ValidateToken(token)
		if err != nil {
			// Токен невалиден - продолжаем без аутентификации
			c.Next()
			return
		}

		// 4. Токен валиден - добавляем в контекст
		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Set("token", token)

		// 5. Продолжаем выполнение
		c.Next()
	}
}
