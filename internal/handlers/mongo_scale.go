package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/installer"
	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/repository"
)

type MongoScaleRequest struct {
	Domain    string `json:"domain"`
	Project   string `json:"project"`
	Cluster   string `json:"cluster"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas  int    `json:"replicas"`
	Storage   string `json:"storage"`
}

func (h *KubeDBHandler) ScaleMongoDB(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MongoScaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
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

	// 1️⃣ Scale in cluster
	if err := installer.ScaleMongo(
		req.Name,
		req.Namespace,
		req.Replicas,
		req.Storage,
	); err != nil {
		http.Error(w, "Cluster scale failed: "+err.Error(), 500)
		return
	}

	// 2️⃣ Update DB
	if err := repository.UpdateMongoScale(
		req.Domain,
		req.Project,
		req.Name,
		req.Replicas,
		req.Storage,
	); err != nil {
		http.Error(w, "DB update failed: "+err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "MongoDB scaled successfully",
	})
}
