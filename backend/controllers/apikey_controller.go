package controllers

import (
	"context"
	"strings"
	"time"

	"api-platform/config"
	"api-platform/middleware"
	"api-platform/models"
	"api-platform/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateKeyRequest defines the input payload for creating a new API key
type CreateKeyRequest struct {
	Name string `json:"name"`
}

// CreateKeyResponse wraps the created API key model with the one-time raw secret
type CreateKeyResponse struct {
	models.APIKey
	RawKey  string `json:"raw_key"`
	Warning string `json:"warning"`
}

// CreateKey generates a new cryptographically secure API key for the authenticated user
func CreateKey(c *fiber.Ctx) error {
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

	var req CreateKeyRequest
	// Body is optional, default name applied if missing or empty
	_ = c.BodyParser(&req)

	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = "Default API Key"
	}

	// Generate secure raw key (sk_live_...)
	rawKey, err := utils.GenerateAPIKey()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate cryptographically secure API key",
		})
	}

	// Compute SHA-256 hash and masked representation
	keyHash := utils.HashAPIKey(rawKey)
	maskedKey := utils.MaskAPIKey(rawKey)

	apiKeyDoc := models.APIKey{
		ID:        primitive.NewObjectID(),
		UserID:    userObjectID,
		KeyHash:   keyHash,
		MaskedKey: maskedKey,
		Name:      name,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	keysCollection := config.GetCollection("apikeys")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = keysCollection.InsertOne(ctx, apiKeyDoc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to persist API key to database",
		})
	}

	// Return the raw key ONLY once
	return c.Status(fiber.StatusCreated).JSON(CreateKeyResponse{
		APIKey:  apiKeyDoc,
		RawKey:  rawKey,
		Warning: "Store this API key safely! It will NOT be shown again and cannot be retrieved later.",
	})
}

// ListKeys returns all API keys belonging to the authenticated user
func ListKeys(c *fiber.Ctx) error {
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

	keysCollection := config.GetCollection("apikeys")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userObjectID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := keysCollection.Find(ctx, filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query API keys from database",
		})
	}
	defer cursor.Close(ctx)

	keys := make([]models.APIKey, 0)
	if err := cursor.All(ctx, &keys); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode API keys list",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"keys":  keys,
		"total": len(keys),
	})
}

// RevokeKey deactivates or deletes an API key, ensuring ownership by the authenticated user
func RevokeKey(c *fiber.Ctx) error {
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

	keyIDHex := c.Params("id")
	keyObjectID, err := primitive.ObjectIDFromHex(keyIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid API key ID parameter",
		})
	}

	keysCollection := config.GetCollection("apikeys")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if hard delete was explicitly requested via ?mode=delete
	mode := strings.ToLower(c.Query("mode"))
	if mode == "delete" {
		res, err := keysCollection.DeleteOne(ctx, bson.M{
			"_id":     keyObjectID,
			"user_id": userObjectID,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete API key record",
			})
		}
		if res.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "API key not found or does not belong to you",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "API key permanently deleted",
			"id":      keyIDHex,
		})
	}

	// Default behavior: Deactivate key (sets is_active = false)
	res, err := keysCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":     keyObjectID,
			"user_id": userObjectID,
		},
		bson.M{
			"$set": bson.M{
				"is_active": false,
			},
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update API key status",
		})
	}

	if res.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "API key not found or does not belong to you",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "API key successfully revoked",
		"id":        keyIDHex,
		"is_active": false,
	})
}
