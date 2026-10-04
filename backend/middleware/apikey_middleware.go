package middleware

import (
	"context"
	"strings"
	"time"

	"api-platform/config"
	"api-platform/models"
	"api-platform/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ValidateAPIKey verifies the x-api-key header against hashed MongoDB records and attaches client identity to context
func ValidateAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawKey := strings.TrimSpace(c.Get("x-api-key"))
		if rawKey == "" {
			// Check standard casing fallback
			rawKey = strings.TrimSpace(c.Get("X-API-Key"))
		}

		if rawKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing required API key. Please provide 'x-api-key' in request headers.",
			})
		}

		// Compute SHA-256 hash of the incoming key
		keyHash := utils.HashAPIKey(rawKey)

		keysCollection := config.GetCollection("apikeys")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var apiKeyDoc models.APIKey
		err := keysCollection.FindOne(ctx, bson.M{
			"key_hash":  keyHash,
			"is_active": true,
		}).Decode(&apiKeyDoc)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid or revoked API key",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error while validating API key",
			})
		}

		// Fetch the associated user to determine their tier (Free vs Premium)
		usersCollection := config.GetCollection("users")
		var userDoc models.User
		err = usersCollection.FindOne(ctx, bson.M{
			"_id": apiKeyDoc.UserID,
		}).Decode(&userDoc)

		tier := models.TierFree
		if err == nil && userDoc.Tier != "" {
			tier = userDoc.Tier
		}

		// Store identity properties in Fiber Locals for downstream handlers
		c.Locals("apiKeyID", apiKeyDoc.ID)
		c.Locals("apiKeyIDHex", apiKeyDoc.ID.Hex())
		c.Locals("userID", apiKeyDoc.UserID)
		c.Locals("userIDHex", apiKeyDoc.UserID.Hex())
		c.Locals("tier", tier)

		return c.Next()
	}
}

// GetContextAPIKeyID retrieves the primitive.ObjectID of the active API key
func GetContextAPIKeyID(c *fiber.Ctx) (primitive.ObjectID, bool) {
	if val := c.Locals("apiKeyID"); val != nil {
		if id, ok := val.(primitive.ObjectID); ok {
			return id, true
		}
	}
	return primitive.NilObjectID, false
}
