package router

import (
	"log"

	"tabler-api/config"
	"tabler-api/internal/handler"
	"tabler-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Inicializa o middleware de autenticação Auth0
	authMiddleware, err := middleware.NewAuth0Middleware(&cfg.Auth0)
	if err != nil {
		log.Fatalf("Failed to initialize Auth0 middleware: %v", err)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User routes (não protegidas)
		users := v1.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Rotas protegidas que requerem autenticação
		protected := v1.Group("/")
		protected.Use(authMiddleware.Authenticate())
		{
			// Profile route - retorna dados do usuário autenticado
			profileHandler := handler.NewProfileHandler()
			protected.GET("/profile", profileHandler.GetProfile)
		}
	}

	return router
} 