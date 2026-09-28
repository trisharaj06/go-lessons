package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goose-neon-demo/auth"
	"goose-neon-demo/config"
)

func RequireAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(auth.AccessTokenCookieName)
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		userID, err := auth.ParseAccessToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
