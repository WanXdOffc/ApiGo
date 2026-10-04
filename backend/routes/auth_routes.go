package routes

import (
	"api-platform/controllers"
	"api-platform/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes sets up all authentication endpoints under /api/auth
func SetupAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	// Manual authentication (Email/Password)
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	// Protected identity verification route to validate JWT tokens
	auth.Get("/me", middleware.Protected(), func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "authenticated",
			"user_id": middleware.GetUserID(c),
			"tier":    middleware.GetUserTier(c),
			"email":   c.Locals("email"),
		})
	})

	// OAuth Placeholders (Explicitly deferred until the very end)
	auth.Get("/google", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error":   "Google OAuth authentication is deferred and not implemented yet",
			"code":    fiber.StatusNotImplemented,
			"message": "Please use manual credentials (/api/auth/register or /api/auth/login)",
		})
	})

	auth.Get("/github", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error":   "GitHub OAuth authentication is deferred and not implemented yet",
			"code":    fiber.StatusNotImplemented,
			"message": "Please use manual credentials (/api/auth/register or /api/auth/login)",
		})
	})
}
