package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/metallb"
	"dbaas-orcastrator/internal/models"
	"dbaas-orcastrator/internal/repository"
)

type MetalLBRequest struct {
	Domain    string   `json:"domain"`
	Project   string   `json:"project"`
	Cluster   string   `json:"cluster"`
	Name      string   `json:"name"`
	Addresses []string `json:"addresses"`
}

func (h *KubeDBHandler) AddMetalLBPool(w http.ResponseWriter, r *http.Request) {

	var req MetalLBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	token, err := getBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 1️⃣ kubeconfig
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

	// 2️⃣ Generate YAML
	yamlPath := "./scripts/metallb-pool.yaml"
	if err := metallb.GenerateYAML(yamlPath, metallb.TemplateData{
		PoolName:          req.Name,
		Addresses:         req.Addresses,
		AdvertisementName: req.Name, // SAME NAME ✔
	}); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 3️⃣ Apply
	if err := metallb.ApplyYAML(
		h.Cfg.KubeconfigPath, // 🔥 tenant kubeconfig
		yamlPath,
	); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 4️⃣ Save DB
	err = repository.SaveMetalLBPool(models.MetalLBPool{
		Domain:            req.Domain,
		Project:           req.Project,
		Cluster:           req.Cluster,
		PoolName:          req.Name,
		Addresses:         req.Addresses,
		AdvertisementName: req.Name,
	})
	if err != nil {
		http.Error(w, "Mongo save failed", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "MetalLB pool created successfully",
	})
}
