package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"dbaas-orcastrator/internal/installer"
	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/models"
	"dbaas-orcastrator/internal/repository"
)

type RedisCreateRequest struct {
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

func (h *KubeDBHandler) RedisCreateRequest(w http.ResponseWriter, r *http.Request) {

	var req RedisCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	record := models.DatabaseRecord{
		Domain:      req.Domain,
		Project:     req.Project,
		Cluster:     req.Cluster,
		DBType:      "redis",
		Name:        req.Name,
		Namespace:   req.Namespace,
		Version:     req.Version,
		Username:    req.Username,
		Password:    req.Password,
		Replicas:    req.Replicas,
		Storage:     req.Storage,
		MetalLBPool: req.MetalLBPool,
		Status:      "CREATING",
		CreatedAt:   time.Now(),
	}

	token, err := getBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 1️⃣ Download kubeconfig
	if err := kubeconfig.Download(
		h.Cfg.CompassBaseURL,
		token,
		req.Domain,
		req.Project,
		req.Cluster,
		h.Cfg.KubeconfigPath,
	); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// 🔥 DB INSERT FIRST
	if err := repository.InsertMongoDB(record); err != nil {
		http.Error(w, "Mongo insert failed", 500)
		return
	}

	// 🔥 RUN INSTALL ASYNC
	go installer.InstallRedis(record)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "Redis creation started",
	})
}
