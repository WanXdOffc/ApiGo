package main

import (
	"log"
	"os"

	"api-platform/config"
	"api-platform/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: .env file not found, using system environment variables")
	}

	// Initialize MongoDB connection before starting server
	if _, err := config.ConnectDB(); err != nil {
		log.Fatalf("Database connection initialization failed: %v", err)
	}
	defer config.DisconnectDB()

	// Initialize Redis connection for rate limiting (with in-memory fallback if unavailable)
	if _, err := config.ConnectRedis(); err != nil {
		log.Printf("Notice: Redis initialization (%v). Sliding-window rate limiter will operate with in-memory fallback.\n", err)
	}
	defer config.CloseRedis()

	// Resolve server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := fiber.New(fiber.Config{
		AppName: "B2D API Platform SaaS - Core API",
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://127.0.0.1:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))

	// API Routes
	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
		})
	})

	// Setup Authentication Routes (/api/auth)
	routes.SetupAuthRoutes(api)

	// Setup API Key Lifecycle Routes (/api/keys)
	routes.SetupAPIKeyRoutes(api)

	// Setup Public Developer API Routes (/api/v1) with API Key validation, Rate Limiting, and Audit Logging
	routes.SetupPublicAPIRoutes(api)

	log.Printf("Starting Core API Server on port :%s...\n", port)
	log.Fatal(app.Listen(":" + port))
}
