package main

import (
	"calendar-backend/internal/database"
	"calendar-backend/internal/handlers"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./calendar.db"
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// API routes
	api := e.Group("/api")
	// Health check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]bool{"ok": true})
	})
	api.GET("/events", handlers.GetEvents)
	api.GET("/events/:id", handlers.GetEvent)
	api.POST("/events", handlers.CreateEvent)
	api.PUT("/events/:id", handlers.UpdateEvent)
	api.DELETE("/events/:id", handlers.DeleteEvent)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "18082"
	}
	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
