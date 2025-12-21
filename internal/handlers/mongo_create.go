package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"dbaas-orcastrator/internal/installer"
	"dbaas-orcastrator/internal/models"
	"dbaas-orcastrator/internal/repository"
)

type MongoCreateRequest struct {
	Domain  string `json:"domain"`
	Project string `json:"project"`
	Cluster string `json:"cluster"`

	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Version   string `json:"version"`

	Username string `json:"username"`
	Password string `json:"password"`

	Replicas    int    `json:"replicas"`
	ReplicaSet  string `json:"replicaSet"`
	Storage     string `json:"storage"`
	MetalLBPool string `json:"metallbPool"`
}

func (h *KubeDBHandler) CreateMongoDB(w http.ResponseWriter, r *http.Request) {

	var req MongoCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	record := models.DatabaseRecord{
		Domain:      req.Domain,
		Project:     req.Project,
		Cluster:     req.Cluster,
		DBType:      "mongo",
		Name:        req.Name,
		Namespace:   req.Namespace,
		Version:     req.Version,
		Username:    req.Username,
		Password:    req.Password,
		Replicas:    req.Replicas,
		ReplicaSet:  req.ReplicaSet,
		Storage:     req.Storage,
		MetalLBPool: req.MetalLBPool,
		Status:      "CREATING",
		CreatedAt:   time.Now(),
	}

	// 🔥 DB INSERT FIRST
	if err := repository.InsertMongoDB(record); err != nil {
		http.Error(w, "Mongo insert failed", 500)
		return
	}

	// 🔥 RUN INSTALL ASYNC
	go installer.InstallMongo(record)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "MongoDB creation started",
	})
}
