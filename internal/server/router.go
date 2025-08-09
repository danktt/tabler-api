package server

import (
	"log/slog"
	"net/http"
	"time"

	"tabler-api/internal/auth"
	"tabler-api/internal/handlers"
	imw "tabler-api/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(logger *slog.Logger) http.Handler {
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()

    r.Use(gin.Recovery())
    r.Use(imw.Logging(logger))

    corsConfig := cors.Config{
        AllowAllOrigins: true,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }
    r.Use(cors.New(corsConfig))

    r.GET("/health", handlers.Health)

    api := r.Group("/api")
    jwksURL := "http://localhost:3000/api/auth/jwks"
    apiAuth := auth.NewAuthMiddleware(jwksURL, logger)

    protected := api.Group("/")
    protected.Use(apiAuth)
    protected.GET("/auth/verify", handlers.VerifyAuth)
    protected.GET("/me", handlers.VerifyAuth)
    protected.GET("/users", handlers.ListUsers)

    return r
}
