package controllers

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

// RegisterRequest defines the input payload for user registration
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest defines the input payload for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse defines the successful authentication response payload
type AuthResponse struct {
	Token     string      `json:"token"`
	TokenType string      `json:"token_type"`
	User      models.User `json:"user"`
}

// Register handles new developer account creation
func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body format",
		})
	}

	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	// Basic validation
	if name == "" || email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email, and password are required fields",
		})
	}

	if len(password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters long",
		})
	}

	usersCollection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user with this email already exists
	var existingUser models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "An account with this email address already exists",
		})
	} else if err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database lookup failed while checking email",
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to securely hash password",
		})
	}

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		Tier:      models.TierFree,
		CreatedAt: time.Now(),
	}

	_, err = usersCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user record in database",
		})
	}

	// Generate JWT for direct authentication upon registration
	token, err := utils.GenerateToken(newUser.ID.Hex(), newUser.Tier, newUser.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User registered but failed to generate session token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Token:     token,
		TokenType: "Bearer",
		User:      newUser,
	})
}

// Login verifies credentials and issues a signed JWT
func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body format",
		})
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	if email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required fields",
		})
	}

	usersCollection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database lookup failed during authentication",
		})
	}

	// Verify bcrypt password hash
	if !utils.CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate JWT token containing User ID and Tier
	token, err := utils.GenerateToken(user.ID.Hex(), user.Tier, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(AuthResponse{
		Token:     token,
		TokenType: "Bearer",
		User:      user,
	})
}
