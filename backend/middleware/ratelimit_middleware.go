package middleware

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"api-platform/config"
	"api-platform/models"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

const (
	// FreeTierLimit is the requests allowed per minute for Free users
	FreeTierLimit = 10
	// PremiumTierLimit is the requests allowed per minute for Premium users
	PremiumTierLimit = 1000
	// RateLimitWindow is the sliding window duration
	RateLimitWindow = 60 * time.Second
)

// inMemoryStore provides a fallback in case Redis is not configured or temporarily unreachable
type inMemoryStore struct {
	sync.Mutex
	records map[string][]int64
}

var memStore = &inMemoryStore{
	records: make(map[string][]int64),
}

func (s *inMemoryStore) isAllowed(key string, limit int) (bool, int) {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixMilli()
	windowStart := now - RateLimitWindow.Milliseconds()

	timestamps := s.records[key]
	valid := make([]int64, 0, len(timestamps))
	for _, ts := range timestamps {
		if ts > windowStart {
			valid = append(valid, ts)
		}
	}

	count := len(valid)
	if count >= limit {
		s.records[key] = valid
		return false, 0
	}

	valid = append(valid, now)
	s.records[key] = valid
	return true, limit - count - 1
}

// RateLimit enforces sliding-window request limits based on user ID and tier (Free vs Premium)
func RateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Identify client by UserID, with fallback to Client IP
		userID := ""
		if val := c.Locals("userIDHex"); val != nil {
			if idStr, ok := val.(string); ok && idStr != "" {
				userID = idStr
			}
		}
		if userID == "" {
			userID = c.IP()
		}

		// Determine user tier and applicable rate limit
		tier := models.TierFree
		if val := c.Locals("tier"); val != nil {
			if t, ok := val.(string); ok && t != "" {
				tier = t
			}
		}

		limit := FreeTierLimit
		if strings.EqualFold(tier, models.TierPremium) {
			limit = PremiumTierLimit
		}

		// 1. If Redis client is active, use atomic sliding window with Redis Sorted Sets
		if config.RedisClient != nil {
			redisKey := fmt.Sprintf("ratelimit:%s", userID)
			now := time.Now().UnixMilli()
			windowStart := now - RateLimitWindow.Milliseconds()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			pipe := config.RedisClient.TxPipeline()
			// Remove expired timestamps
			pipe.ZRemRangeByScore(ctx, redisKey, "-inf", strconv.FormatInt(windowStart, 10))
			// Count hits in current window
			countCmd := pipe.ZCard(ctx, redisKey)
			// Record current request with a unique member key
			uniqueMember := fmt.Sprintf("%d:%d", now, rand.Int63())
			pipe.ZAdd(ctx, redisKey, redis.Z{
				Score:  float64(now),
				Member: uniqueMember,
			})
			// Expire set after window duration + buffer
			pipe.Expire(ctx, redisKey, RateLimitWindow+10*time.Second)

			if _, err := pipe.Exec(ctx); err != nil {
				log.Printf("Redis rate limit error (%v), falling back to in-memory check\n", err)
			} else {
				currentCount := int(countCmd.Val())
				remaining := limit - currentCount - 1
				if remaining < 0 {
					remaining = 0
				}

				c.Set("X-RateLimit-Limit", strconv.Itoa(limit))
				c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

				if currentCount >= limit {
					c.Set("Retry-After", "60")
					return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
						"error":               "Too Many Requests",
						"message":             fmt.Sprintf("Rate limit exceeded for %s tier (%d requests/min). Upgrade to Premium for higher limits.", tier, limit),
						"tier":                tier,
						"limit":               limit,
						"retry_after_seconds": 60,
					})
				}

				return c.Next()
			}
		}

		// 2. In-memory sliding-window fallback
		allowed, remaining := memStore.isAllowed(userID, limit)
		c.Set("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			c.Set("Retry-After", "60")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":               "Too Many Requests",
				"message":             fmt.Sprintf("Rate limit exceeded for %s tier (%d requests/min). Upgrade to Premium for higher limits.", tier, limit),
				"tier":                tier,
				"limit":               limit,
				"retry_after_seconds": 60,
			})
		}

		return c.Next()
	}
}
