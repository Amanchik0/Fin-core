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
				"error":   "Authorization header is required",
			})
			c.Abort()
			return
		}

		user, err := authClient.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired token",
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
		origin := c.Request.Header.Get("Origin")

		// Список разрешенных origins
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8081",  // Expo web dev server
			"http://localhost:8082",  // Альтернативный порт Expo
			"http://localhost:19006", // Expo web альтернативный порт
			"http://localhost:19000", // Metro bundler
		}

		// Проверяем, разрешен ли origin
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cookie")
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
