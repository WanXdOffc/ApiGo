package routes

import (
	"api-platform/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupPublicAPIRoutes configures public developer API endpoints under /api/v1
func SetupPublicAPIRoutes(router fiber.Router) {
	v1 := router.Group("/v1")

	// Apply Security, Rate Limiting, and Telemetry middlewares:
	// 1. AuditLogMiddleware: records request method, path, latency, and status asynchronously
	// 2. ValidateAPIKey: verifies x-api-key, verifies hash against MongoDB, injects User & Tier
	// 3. RateLimit: enforces sliding-window quotas via Redis based on User tier
	v1.Use(middleware.AuditLogMiddleware())
	v1.Use(middleware.ValidateAPIKey())
	v1.Use(middleware.RateLimit())

	// Mock public data endpoint protected by API Key
	v1.Get("/data", func(c *fiber.Ctx) error {
		tier := c.Locals("tier")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data":    "Here is your protected data!",
			"tier":    tier,
			"message": "API key successfully authenticated and request logged.",
		})
	})
}
