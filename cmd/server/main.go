package main

import (
	"log"
	"net/http"

	"dbaas-orcastrator/internal/config"
	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/handlers"
	"dbaas-orcastrator/internal/routes"
)

func main() {
	cfg := config.Load()

	db.Connect(cfg)

	kubeHandler := &handlers.KubeDBHandler{Cfg: cfg}
	routes.Register(kubeHandler)

	log.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
