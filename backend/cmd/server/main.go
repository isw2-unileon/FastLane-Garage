// Package main is the entry point for the backend server.
// It initializes the database, sets up the HTTP router and starts the server
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/config"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/handlers"
)

// logger is a structured logger instance used throughout the application
// It writes JSON-formatted logs to satndard output for easy parsing and monitoring
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// Is the application entry point
// It perfoms the following operations in order:
// 1. Load configuration form environment variables
// 2. Initialize the SQLite database
// 3. Run database migrations
// 4. Seed the database with sample data
// 5. Configure and start the HTTP server
// 6. Handle graceful shutdown on interrupt signals
func main() {
	// Create a background context for use during initilization
	ctx := context.Background()

	// Load configuration form the environment variables
	// The config package provides defaults if variables are not set
	cfg := config.Load()

	// Initialize the SQLite database connection
	// This creates a connection to "garage" in the current directory
	db, err := database.Init("garage.db")
	if err != nil {
		logger.Error("database initialization failed", "error", err)
		os.Exit(1)
	}

	// Run database migrations to create or update the schema
	// This ensures the "parts" table exists with the correct structure
	if err := database.Migrate(db); err != nil {
		logger.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	// Seed the database with sample parts data
	// This is idempotent-if data already exists, it will be skiped
	if err := database.SeedParts(db); err != nil {
		logger.Error("database seeding failed", "error", err)
		os.Exit(1)
	}

	// Set Gin's mode based on the configuration
	// In development, mode is "debug" for verbose logging
	// In production, mode is "release" for minimal output
	gin.SetMode(cfg.GinMode)

	// Create a new Gin router instance
	r := gin.New()

	// Attach middleware to log HTTP requests and recover from panics
	r.Use(gin.Logger(), gin.Recovery())

	// Register the health check endpoint
	// This endpoint is used by load balancers and monitoring systems
	// to verify that the server is running and healthy
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create an API route group with the "/api" prefix
	// All routes defined within this group will be accessible at "/api/*"
	api := r.Group("/api")
	{
		// Register the GET /api/parts endpoint
		// This endpoint retrieves all automobile parts from the database
		// for the frontend to display in the interactive SVG
		api.GET("/parts", handlers.GetParts(db))

		// Exist test endpoint
		// Register the GET /api/hello endpoint
		// This is a simple endpoint for testing the API is working
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello from the API"})
		})
	}

	// Create an HTTP serverinstance with the configured router
	// Timeout are set to prevent slow clients from consuming resources indefinitely
	srv := &http.Server{
		Addr:         ":" + cfg.Port, // Liten on all interfaces on the configured port
		Handler:      r,              // Use the Gin router as the HTTP handler
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create a context that can be cancelled by interrupt or termination signals
	// This allows gracefull shutdown when the process receives SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the HTTP server in a separate goroutine so it runs concurrently
	// with the signal handling code below
	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal (Ctrl + C)
	// This block until the context is cancelled by a signal
	<-ctx.Done()
	slog.Info("shutting down server")

	// Create a new context with a 5-second timeout for graceful shutdown
	// This ensures all in-flight requests complete before the server stop,
	// but the server won't forever if clients are slow
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shut down the server
	// This stops accepting new connections and wait for existing ones to close
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}

	// Log the final message indicating the server has stopped
	logger.Info("server stopped")
}
