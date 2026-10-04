package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLog tracks every inbound request executed using an API key
type AuditLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	APIKeyID   primitive.ObjectID `bson:"api_key_id" json:"api_key_id"`
	Endpoint   string             `bson:"endpoint" json:"endpoint"`
	Method     string             `bson:"method" json:"method"`
	StatusCode int                `bson:"status_code" json:"status_code"`
	Latency    int64              `bson:"latency" json:"latency"` // Latency in milliseconds
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
}
