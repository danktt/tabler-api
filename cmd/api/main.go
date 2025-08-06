package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tabler-api/config"
	"tabler-api/internal/handler"
	"tabler-api/internal/repository"
	"tabler-api/internal/router"
	"tabler-api/internal/service"
	"tabler-api/pkg/db"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := router.SetupRouter(userHandler, cfg)

	// Create HTTP server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
} 