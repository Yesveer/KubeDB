package services

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FixMongoURI
// - URL encode username/password
// - force admin db
// - add authSource=admin
// - add directConnection=true
func FixMongoURI(uri string) (string, error) {

	// log.Println("🔧 Original Mongo URI:", uri)

	// Split scheme
	parts := strings.SplitN(uri, "://", 2)
	if len(parts) != 2 {
		return uri, nil
	}

	scheme := parts[0]
	rest := parts[1]

	// Split auth and host
	if !strings.Contains(rest, "@") {
		// No auth present
		final := scheme + "://" + rest + "/admin?authSource=admin&directConnection=true"
		log.Println("🔧 Fixed Mongo URI:", final)
		return final, nil
	}

	authHost := strings.SplitN(rest, "@", 2)
	auth := authHost[0] // username:password
	host := authHost[1]

	userPass := strings.SplitN(auth, ":", 2)
	if len(userPass) != 2 {
		return uri, nil
	}

	user := url.QueryEscape(userPass[0])
	pass := url.QueryEscape(userPass[1])

	// Remove any existing db/query
	host = strings.Split(host, "/")[0]

	finalURI := scheme + "://" +
		user + ":" + pass + "@" +
		host +
		"/admin?authSource=admin&directConnection=true"

	// log.Println("🔧 Fixed Mongo URI:", finalURI)

	return finalURI, nil
}

func NewMongoClient(uri string) (*mongo.Client, error) {

	fixedURI, err := FixMongoURI(uri)
	if err != nil {
		return nil, err
	}

	// log.Println("🚀 Connecting to MongoDB...")
	// log.Println("🔗 Mongo URI:", fixedURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fixedURI))
	if err != nil {
		log.Println("❌ Mongo connection failed:", err)
		return nil, err
	}

	// 🔥 Ping to ensure connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Println("❌ Mongo ping failed:", err)
		return nil, err
	}

	// log.Println("✅ MongoDB connected successfully")
	return client, nil
}
