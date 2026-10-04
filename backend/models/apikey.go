package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIKey represents an API key credential generated for developer authentication
type APIKey struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	KeyHash   string             `bson:"key_hash" json:"-"`
	MaskedKey string             `bson:"masked_key" json:"masked_key"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
}
