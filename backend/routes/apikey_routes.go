package routes

import (
	"api-platform/controllers"
	"api-platform/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAPIKeyRoutes sets up all API key management endpoints under /api/keys
func SetupAPIKeyRoutes(router fiber.Router) {
	// Secure the entire /api/keys group with JWT authentication middleware
	keys := router.Group("/keys", middleware.Protected())

	// Connect API key CRUD endpoints
	keys.Post("/", controllers.CreateKey)
	keys.Get("/", controllers.ListKeys)
	keys.Delete("/:id", controllers.RevokeKey)
}
