package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"goose-neon-demo/config"
	"goose-neon-demo/db"
	"goose-neon-demo/handlers"
	"goose-neon-demo/mailer"
	"goose-neon-demo/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on process env vars")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.InitDB(cfg.DBConnString); err != nil {
		log.Fatal(err)
	}

	m := mailer.New(cfg)

	router := gin.Default()

	authGroup := router.Group("/auth")
	authGroup.POST("/magic-link", handlers.MagicLink(cfg, m))
	authGroup.GET("/verify", handlers.Verify(cfg))
	authGroup.POST("/exchange", handlers.Exchange(cfg))
	authGroup.POST("/logout", handlers.Logout(cfg))
	authGroup.POST("/refresh", handlers.Refresh(cfg))

	usersGroup := router.Group("/users")
	usersGroup.Use(middleware.RequireAuth(cfg))
	usersGroup.GET("/me", handlers.Me)

	router.Run(":" + cfg.Port)
}
