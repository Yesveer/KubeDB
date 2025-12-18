package db

import (
	"context"
	"log"
	"time"

	"dbaas-orcastrator/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("❌ Mongo connect failed:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Mongo ping failed:", err)
	}

	Client = client
	log.Println("✅ MongoDB connected")
}
