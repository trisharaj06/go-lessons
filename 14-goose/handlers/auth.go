package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"goose-neon-demo/auth"
	"goose-neon-demo/config"
	"goose-neon-demo/mailer"
	"goose-neon-demo/models"
)

type magicLinkRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func MagicLink(cfg *config.Config, m *mailer.Mailer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req magicLinkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "a valid email is required"})
			return
		}

		user, err := models.FindOrCreateUserByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not process request"})
			return
		}

		rawToken, err := auth.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not process request"})
			return
		}

		expiresAt := time.Now().Add(cfg.MagicLinkTTL)
		if err := models.CreateMagicLinkToken(user.ID, auth.HashToken(rawToken), expiresAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not process request"})
			return
		}

		verifyURL := fmt.Sprintf("%s/auth/verify?token=%s", cfg.AppBaseURL, rawToken)
		if err := m.SendMagicLink(user.Email, verifyURL); err != nil {
			log.Printf("failed to send magic link email: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "if an account exists for that email, a magic link has been sent"})
	}
}

func Verify(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"valid": false})
			return
		}

		_, err := models.FindValidMagicLinkTokenByHash(auth.HashToken(token))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"valid": false})
			return
		}

		c.JSON(http.StatusOK, gin.H{"valid": true})
	}
}

type exchangeRequest struct {
	Token string `json:"token" binding:"required"`
}

func Exchange(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req exchangeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "token is required"})
			return
		}

		linkToken, err := models.ConsumeMagicLinkToken(auth.HashToken(req.Token))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusBadRequest, gin.H{"message": "invalid or expired token"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not process request"})
			return
		}

		user, err := models.FindUserByID(linkToken.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not process request"})
			return
		}

		if err := issueSession(c, cfg, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not process request"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email, "created_at": user.CreatedAt})
	}
}

func Logout(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if refreshToken, err := c.Cookie(auth.RefreshTokenCookieName); err == nil && refreshToken != "" {
			_ = models.DeleteSessionByRefreshHash(auth.HashToken(refreshToken))
		}

		auth.ClearAuthCookies(c, cfg)
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}

func Refresh(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie(auth.RefreshTokenCookieName)
		if err != nil || refreshToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		session, err := models.FindSessionByRefreshHash(auth.HashToken(refreshToken))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		user, err := models.FindUserByID(session.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not process request"})
			return
		}

		if err := models.DeleteSession(session.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not process request"})
			return
		}

		if err := issueSession(c, cfg, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not process request"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "refreshed"})
	}
}

// issueSession creates a new session row and sets fresh access/refresh cookies.
func issueSession(c *gin.Context, cfg *config.Config, user *models.User) error {
	rawRefreshToken, err := auth.GenerateToken()
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(cfg.RefreshTokenTTL)
	if _, err := models.CreateSession(user.ID, auth.HashToken(rawRefreshToken), expiresAt); err != nil {
		return err
	}

	accessToken, err := auth.IssueAccessToken(user.ID, cfg.JWTSecret, cfg.AccessTokenTTL)
	if err != nil {
		return err
	}

	auth.SetAuthCookies(c, accessToken, rawRefreshToken, cfg)
	return nil
}
