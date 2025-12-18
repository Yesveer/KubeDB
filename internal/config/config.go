package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CompassBaseURL string
	KubeconfigPath string
	MongoURI       string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		CompassBaseURL: os.Getenv("COMPASS_BASE_URL"),
		KubeconfigPath: os.Getenv("KUBECONFIG_PATH"),
		MongoURI:       os.Getenv("MONGO_URI"),
	}

	if cfg.CompassBaseURL == "" || cfg.KubeconfigPath == "" {
		log.Fatal("❌ required env variables missing")
	}

	return cfg
}
