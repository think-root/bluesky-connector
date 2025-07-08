package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/think-root/bluesky-connector/internal/client"
	"github.com/think-root/bluesky-connector/internal/config"
	"github.com/think-root/bluesky-connector/internal/handlers"
	"github.com/think-root/bluesky-connector/internal/logger"
	"github.com/think-root/bluesky-connector/internal/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.Log.Level)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Fatalf("Configuration validation failed: %v", err)
	}

	// Initialize Bluesky client
	blueSkyClient := client.NewBlueSkyClient(cfg)

	// Test authentication
	if err := blueSkyClient.Authenticate(); err != nil {
		logger.Fatalf("Failed to authenticate with Bluesky: %v", err)
	}

	// Set Gin mode
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Initialize handlers
	postHandler := handlers.NewPostHandler(blueSkyClient)

	// Health check route (no authentication required)
	router.GET("/bluesky/api/health", postHandler.HealthCheck)

	// API routes (authenticated routes under /bluesky/api prefix)
	api := router.Group("/bluesky/api")
	api.Use(middleware.APIKeyMiddleware(cfg))
	{
		api.POST("/posts/create", postHandler.CreatePost)
		api.POST("/test/posts/create", postHandler.CreateTestPost)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting Bluesky Connector server on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}