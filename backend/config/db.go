package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// Client is the global MongoDB client instance
	Client *mongo.Client
	// DB is the global MongoDB database instance
	DB *mongo.Database
)

// ConnectDB establishes a connection to MongoDB, validates via ping, and assigns global client & DB
func ConnectDB() (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable is not defined")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "portfolio"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	// Ping primary to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	Client = client
	DB = client.Database(dbName)

	log.Printf("Successfully connected and pinged MongoDB database: [%s]\n", dbName)
	return Client, nil
}

// GetCollection returns a handle to a MongoDB collection in the default database
func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Fatalf("Database connection has not been initialized yet")
	}
	return DB.Collection(collectionName)
}

// DisconnectDB cleanly terminates the MongoDB client connection
func DisconnectDB() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := Client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v\n", err)
		} else {
			log.Println("MongoDB connection closed cleanly.")
		}
	}
}
