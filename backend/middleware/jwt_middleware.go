package middleware

import (
	"strings"

	"api-platform/utils"

	"github.com/gofiber/fiber/v2"
)

// Protected verifies the Bearer JWT token from the Authorization header
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header. Expected format: Bearer <token>",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Malformed Authorization header. Format must be: Bearer <token>",
			})
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token string cannot be empty",
			})
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid, expired, or malformed authentication token",
			})
		}

		// Store parsed claims in Fiber Locals for downstream handlers
		c.Locals("userID", claims.UserID)
		c.Locals("tier", claims.Tier)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

// GetUserID extracts the authenticated user's ID stored in Fiber Locals
func GetUserID(c *fiber.Ctx) string {
	if val := c.Locals("userID"); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserTier extracts the authenticated user's tier stored in Fiber Locals
func GetUserTier(c *fiber.Ctx) string {
	if val := c.Locals("tier"); val != nil {
		if tier, ok := val.(string); ok {
			return tier
		}
	}
	return ""
}
