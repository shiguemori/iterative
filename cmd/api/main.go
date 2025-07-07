// @title B3 Challenge API
// @version 1.0
// @description API B3.
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "b3-challenge/docs"
	"b3-challenge/internal/api"
	"b3-challenge/internal/config"
	"b3-challenge/internal/repository"
	"b3-challenge/internal/service"
	redis "b3-challenge/pkg/cache"
	"b3-challenge/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		err := database.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}()

	// Initialize dependencies
	tradeRepo := repository.NewTradeRepository(database.GetDB())
	redisClient := redis.NewRedisClient(&cfg.Redis)
	tradeService := service.NewTradeService(tradeRepo, redisClient)

	// Setup router
	router := api.SetupRouter(tradeService)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.API.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting B3 Challenge API server on port %d", cfg.API.Port)
		log.Printf("API Documentation available at: http://localhost:%d/", cfg.API.Port)
		log.Printf("Health check available at: http://localhost:%d/health", cfg.API.Port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
