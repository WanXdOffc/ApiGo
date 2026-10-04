package middleware

import (
	"context"
	"log"
	"time"

	"api-platform/config"
	"api-platform/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLogMiddleware captures the full request lifecycle and asynchronously logs it to MongoDB
func AuditLogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Execute downstream handler first
		err := c.Next()

		// Calculate latency in milliseconds
		latencyMs := time.Since(startTime).Milliseconds()
		statusCode := c.Response().StatusCode()

		// Extract API key ID from context if populated by API key middleware
		var apiKeyID primitive.ObjectID
		if val := c.Locals("apiKeyID"); val != nil {
			if id, ok := val.(primitive.ObjectID); ok {
				apiKeyID = id
			}
		}

		// CRITICAL CONCURRENCY SAFETY:
		// Fiber recycles *fiber.Ctx in a fasthttp memory pool immediately after this handler returns.
		// We MUST extract and copy all string/scalar values into local variables BEFORE launching
		// the goroutine so the background task never references the pooled Fiber context.
		copiedMethod := string(c.Request().Header.Method())
		copiedEndpoint := string(c.Request().URI().Path())
		capturedTimestamp := startTime

		// Spawn asynchronous goroutine to persist audit log without delaying API response
		go func(apiKeyID primitive.ObjectID, endpoint, method string, statusCode int, latencyMs int64, ts time.Time) {
			auditDoc := models.AuditLog{
				ID:         primitive.NewObjectID(),
				APIKeyID:   apiKeyID,
				Endpoint:   endpoint,
				Method:     method,
				StatusCode: statusCode,
				Latency:    latencyMs,
				Timestamp:  ts,
			}

			// Use isolated background context with a timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			coll := config.GetCollection("auditlogs")
			if coll != nil {
				if _, insertErr := coll.InsertOne(ctx, auditDoc); insertErr != nil {
					log.Printf("AuditLog async insertion error: %v\n", insertErr)
				}
			}
		}(apiKeyID, copiedEndpoint, copiedMethod, statusCode, latencyMs, capturedTimestamp)

		return err
	}
}
