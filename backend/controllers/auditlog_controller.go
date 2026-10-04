package controllers

import (
	"context"
	"time"

	"api-platform/config"
	"api-platform/middleware"
	"api-platform/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ListAuditLogs returns all request logs associated with any API key owned by the authenticated user
func ListAuditLogs(c *fiber.Ctx) error {
	userIDHex := middleware.GetUserID(c)
	if userIDHex == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: user ID missing from context",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Retrieve all API keys owned by this user
	keysCollection := config.GetCollection("apikeys")
	keyCursor, err := keysCollection.Find(ctx, bson.M{"user_id": userObjectID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to look up user API keys",
		})
	}
	defer keyCursor.Close(ctx)

	var userKeys []models.APIKey
	if err := keyCursor.All(ctx, &userKeys); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode user API keys",
		})
	}

	// If user has no keys, return empty logs array immediately
	if len(userKeys) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"logs":  []models.AuditLog{},
			"total": 0,
		})
	}

	keyIDs := make([]primitive.ObjectID, len(userKeys))
	for i, k := range userKeys {
		keyIDs[i] = k.ID
	}

	// 2. Query audit logs matching any of user's key IDs
	logsCollection := config.GetCollection("auditlogs")
	filter := bson.M{"api_key_id": bson.M{"$in": keyIDs}}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(100)

	logCursor, err := logsCollection.Find(ctx, filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query audit logs from database",
		})
	}
	defer logCursor.Close(ctx)

	logs := make([]models.AuditLog, 0)
	if err := logCursor.All(ctx, &logs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode audit logs",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"logs":  logs,
		"total": len(logs),
	})
}
