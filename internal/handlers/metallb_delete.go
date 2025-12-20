package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/metallb"
	"dbaas-orcastrator/internal/repository"
)

type DeleteMetalLBRequest struct {
	Domain  string `json:"domain"`
	Project string `json:"project"`
	Cluster string `json:"cluster"`
	Name    string `json:"name"`
}

func (h *KubeDBHandler) DeleteMetalLBPool(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteMetalLBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	token, err := getBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

	// 2️⃣ Delete from cluster
	if err := metallb.DeletePool(
		h.Cfg.KubeconfigPath,
		req.Name,
	); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 3️⃣ Delete from DB
	if err := repository.DeleteMetalLBPool(
		req.Domain,
		req.Project,
		req.Cluster,
		req.Name,
	); err != nil {
		http.Error(w, "Mongo delete failed: "+err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "MetalLB pool deleted successfully",
	})
}
