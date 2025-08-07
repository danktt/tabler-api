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

	// Adiciona middleware de CORS
	router.Use(middleware.CORSMiddleware(&cfg.CORS))

	// Inicializa o middleware de autenticação Better Auth
	var authMiddleware *middleware.BetterAuthMiddleware
	var err error
	
	// Only initialize auth middleware if not in development mode
	if cfg.Env != "development" {
		authMiddleware, err = middleware.NewBetterAuthMiddleware(&cfg.BetterAuth)
		if err != nil {
			log.Fatalf("Failed to initialize Better Auth middleware: %v", err)
		}
	} else {
		log.Println("Running in development mode - auth middleware disabled")
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	// Rota de teste CORS (sem autenticação)
	router.GET("/api/v1/test-cors", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "CORS is working correctly!",
			"origin":  c.GetHeader("Origin"),
			"method":  c.Request.Method,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User routes (protegidas por autenticação)
		users := v1.Group("/users")
		if authMiddleware != nil {
			users.Use(authMiddleware.Authenticate())
		}
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Rotas protegidas que requerem autenticação
		protected := v1.Group("/")
		if authMiddleware != nil {
			protected.Use(authMiddleware.Authenticate())
		}
		{
			// Profile route - retorna dados do usuário autenticado
			profileHandler := handler.NewProfileHandler()
			protected.GET("/profile", profileHandler.GetProfile)
		}
	}

	return router
} 