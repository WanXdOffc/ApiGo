package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TierFree    = "Free"
	TierPremium = "Premium"
)

// User represents a registered developer account on the platform
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Tier      string             `bson:"tier" json:"tier"` // "Free" or "Premium"
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
