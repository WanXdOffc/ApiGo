package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient is the global Redis client instance
var RedisClient *redis.Client

// ConnectRedis parses the REDIS_URL environment variable, verifies connection with a ping, and stores the client
func ConnectRedis() (*redis.Client, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Println("Notice: REDIS_URL not configured. Redis rate limiting will be inactive.")
		return nil, fmt.Errorf("REDIS_URL is not set")
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Warning: Failed to parse REDIS_URL (%v)\n", err)
		return nil, fmt.Errorf("invalid REDIS_URL: %w", err)
	}

	// Set reasonable timeouts for cloud Redis (Upstash)
	opts.DialTimeout = 5 * time.Second
	opts.ReadTimeout = 3 * time.Second
	opts.WriteTimeout = 3 * time.Second

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("Warning: Could not ping Redis at %s: %v\n", opts.Addr, err)
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	RedisClient = client
	log.Printf("Successfully connected and pinged Redis at: [%s]\n", opts.Addr)
	return RedisClient, nil
}

// CloseRedis cleanly terminates the Redis connection pool
func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Error closing Redis client: %v\n", err)
		} else {
			log.Println("Redis client closed cleanly.")
		}
	}
}
