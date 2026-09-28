package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goose-neon-demo/config"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

func SetAuthCookies(c *gin.Context, accessToken, refreshToken string, cfg *config.Config) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AccessTokenCookieName, accessToken, int(cfg.AccessTokenTTL.Seconds()), "/", cfg.CookieDomain, cfg.CookieSecure, true)
	c.SetCookie(RefreshTokenCookieName, refreshToken, int(cfg.RefreshTokenTTL.Seconds()), "/", cfg.CookieDomain, cfg.CookieSecure, true)
}

func ClearAuthCookies(c *gin.Context, cfg *config.Config) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AccessTokenCookieName, "", -1, "/", cfg.CookieDomain, cfg.CookieSecure, true)
	c.SetCookie(RefreshTokenCookieName, "", -1, "/", cfg.CookieDomain, cfg.CookieSecure, true)
}
