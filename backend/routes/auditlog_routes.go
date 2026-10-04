package routes

import (
	"api-platform/controllers"
	"api-platform/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAuditLogRoutes sets up audit log endpoints under /api/logs
func SetupAuditLogRoutes(router fiber.Router) {
	logs := router.Group("/logs", middleware.Protected())

	logs.Get("/", controllers.ListAuditLogs)
}
