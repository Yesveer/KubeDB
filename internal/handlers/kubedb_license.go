package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/installer"
	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/licence"
)

type LicenceRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Domain  string `json:"domain"`
	Project string `json:"project"`
	Cluster string `json:"cluster"`
}

func (h *KubeDBHandler) GenerateLicense(w http.ResponseWriter, r *http.Request) {

	var req LicenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", 400)
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

	// 2️⃣ Get cluster UID
	clusterUID, err := installer.GetClusterUID()
	if err != nil {
		http.Error(w, "Failed to get cluster UID: "+err.Error(), 500)
		return
	}

	// 3️⃣ Call AppsCode licence API
	if err := licence.Generate(req.Name, req.Email, clusterUID); err != nil {
		http.Error(w, "License generation failed: "+err.Error(), 500)
		return
	}

	// 4️⃣ Response
	json.NewEncoder(w).Encode(map[string]string{
		"message": "License email sent successfully",
	})
}
