package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/laundryos/backend/internal/api/handlers"
	"github.com/laundryos/backend/internal/api/middleware"
	"github.com/laundryos/backend/internal/service"
	"github.com/laundryos/backend/pkg/config"
	"github.com/laundryos/backend/pkg/database"
	"github.com/laundryos/backend/pkg/jwt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Connected to database successfully!")

	jwtManager := jwt.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	authService := service.NewAuthService(db, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.TraceID())
	r.Use(middleware.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.Auth(jwtManager), authHandler.Me)
		}
	}

	fmt.Printf("Server starting on port %s\n", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
