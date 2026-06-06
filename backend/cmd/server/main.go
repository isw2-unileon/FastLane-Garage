// Package main is the entry point for the backend server.
// It initializes the database, sets up the HTTP router, and starts the server.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/config"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/handlers"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/repository"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/service"
)

// logger is a structured logger instance used throughout the application.
// It writes JSON-formatted logs to standard output for easy parsing and monitoring.
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// main is the application entry point.
// It performs the following operations in order:
// 1. Load configuration from environment variables
// 2. Initialize the SQLite database
// 3. Run database migrations
// 4. Seed the database with sample data
// 5. Configure and start the HTTP server
// 6. Handle graceful shutdown on interrupt signals.
//
//nolint:funlen // main is naturally long due to dependency injection and route setup
func main() {
	// Create a background context for use during initialization.
	ctx := context.Background()

	// Load configuration from the environment variables.
	// The config package provides defaults if variables are not set.
	cfg := config.Load()

	// Initialize the SQLite database connection.
	// This creates a connection to "garage.db" in the current directory.
	db, err := database.Init("garage.db")
	if err != nil {
		logger.Error("database initialization failed", "error", err)
		os.Exit(1)
	}

	// Run database migrations to create or update the schema.
	// This ensures the "parts" table exists with the correct structure.
	if err := database.Migrate(db); err != nil {
		logger.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	// Seed the database with sample parts data.
	// This is idempotent—if data already exists, it will be skipped.
	if err := database.SeedParts(db); err != nil {
		logger.Error("database seeding failed", "error", err)
		os.Exit(1)
	}

	// Set Gin's mode based on the configuration.
	// In development, mode is "debug" for verbose logging.
	// In production, mode is "release" for minimal output.
	gin.SetMode(cfg.GinMode)

	// Create a new Gin router instance.
	r := gin.New()

	// Attach middleware to log HTTP requests and recover from panics.
	r.Use(gin.Logger(), gin.Recovery())

	// ====================== CORS MIDDLEWARE ======================
	// This allows the React frontend (running on a different port) to access the API.
	corsConfig := cors.DefaultConfig()
	if cfg.CORSAllowOrigin == "*" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = []string{cfg.CORSAllowOrigin}
	}

	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	r.Use(cors.New(corsConfig))

	// Register the health check endpoint.
	// This endpoint is used by load balancers and monitoring systems
	// to verify that the server is running and healthy.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create an API route group with the "/api" prefix.
	// All routes defined within this group will be accessible at "/api/*".
	api := r.Group("/api")

	// Initialize the parts repository with the database connection.
	// The repository is responsible for all database operations related to parts.
	partsRepo := repository.NewPartsRepository(db)

	// Initialize the parts service with the repository.
	// The service contains the business logic for parts operations.
	partsSvc := service.NewPartsService(partsRepo)

	// Initialize the orders repository with the database connection.
	ordersRepo := repository.NewOrdersRepository(db)

	// Initialize the orders service with the repository.
	ordersSvc := service.NewOrdersService(ordersRepo)

	// Initialize the chat repository with the database connection.
	chatRepo := repository.NewChatRepository(db)

	// Initialize the chat service with the repository.
	chatSvc := service.NewChatService(chatRepo, cfg.N8NWebhookURL)

	// Register all CRUD endpoints for parts.
	// Each handler is registered with the HTTP method and route path.

	// ==================== PARTS ENDPOINTS ====================
	// GET /api/parts - Retrieve all parts.
	api.GET("/parts", handlers.GetParts(partsSvc))

	// GET /api/parts/:id - Retrieve a single part by ID.
	api.GET("/parts/:id", handlers.GetPartByID(partsSvc))

	// POST /api/parts - Create a new part.
	api.POST("/parts", handlers.CreatePart(partsSvc))

	// PUT /api/parts/:id - Update an existing part.
	api.PUT("/parts/:id", handlers.UpdatePart(partsSvc))

	// DELETE /api/parts/:id - Delete a part by ID.
	api.DELETE("/parts/:id", handlers.DeletePart(partsSvc))

	// ==================== ORDERS ENDPOINTS ====================
	// GET /api/orders - Retrieve all orders.
	api.GET("/orders", handlers.GetOrders(ordersSvc))

	// GET /api/orders/:id - Retrieve a single order by ID.
	api.GET("/orders/:id", handlers.GetOrderByID(ordersSvc))

	// POST /api/orders- Create a new order.
	api.POST("/orders", handlers.CreateOrder(ordersSvc, partsRepo))

	// PUT /api/orders/:status - Update an order's status.
	api.PUT("/orders/:id/status", handlers.UpdateOrderStatus(ordersSvc))

	// DELETE /api/orders/:id - Delete an order by ID.
	api.DELETE("/orders/:id", handlers.DeleteOrder(ordersSvc))

	// ==================== CHAT ENDPOINTS ====================
	// POST /api/chat/sessions - Create a new session.
	api.POST("/chat/sessions", handlers.CreateChatSession(chatSvc))

	// GET /api/chat/sessions/:id - Get session history.
	api.GET("/chat/sessions/:id", handlers.GetChatSessionHistory(chatSvc))

	// POST /api/chat/sessions/:id/messages - Send a message to n8n.
	api.POST("/chat/sessions/:id/messages", handlers.SendChatMessage(chatSvc))

	// ==================== SAMPLE ENDPOINTS ====================
	// Register the GET /api/hello endpoint.
	// This is a sample endpoint for testing the API is working.
	api.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from the API"})
	})

	// Create an HTTP server instance with the configured router.
	// Timeouts are set to prevent slow clients from consuming resources indefinitely.
	srv := &http.Server{
		Addr:         ":" + cfg.Port, // Listen on all interfaces on the configured port.
		Handler:      r,              // Use the Gin router as the HTTP handler.
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create a context that can be cancelled by interrupt or termination signals.
	// This allows graceful shutdown when the process receives SIGINT or SIGTERM.
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the HTTP server in a separate goroutine so it runs concurrently
	// with the signal handling code below.
	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for a shutdown signal (Ctrl+C or termination signal).
	// This blocks until the context is cancelled by a signal.
	<-ctx.Done()
	slog.Info("shutting down server")

	// Create a new context with a 5-second timeout for graceful shutdown.
	// This ensures all in-flight requests complete before the server stops,
	// but the server won't wait forever if clients are slow.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shut down the server.
	// This stops accepting new connections and waits for existing ones to close.
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}

	// Log the final message indicating the server has stopped.
	logger.Info("server stopped")
}
