package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/installer"
	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/repository"
)

func (h *KubeDBHandler) Check(w http.ResponseWriter, r *http.Request) {

	// 1️⃣ Parse request
	var req InstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", 400)
		return
	}

	// 2️⃣ Bearer token
	token, err := getBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 3️⃣ Download kubeconfig (tenant)
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

	// 4️⃣ Encode kubeconfig
	kubeconfigB64, err := kubeconfig.EncodeFileToBase64(h.Cfg.KubeconfigPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 5️⃣ Check KubeDB via kubectl
	installed, output, err := installer.IsKubeDBInstalled()
	if err != nil {
		http.Error(w, "kubectl failed: "+err.Error(), 500)
		return
	}

	status := "NOT_INSTALLED"
	if installed {
		status = "INSTALLED"

		// Save in DB
		_ = repository.SaveDBaaSRecord(
			req.Domain,
			req.Project,
			req.Cluster,
			status,
			kubeconfigB64,
		)
	}

	// 6️⃣ Response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"installed": status == "INSTALLED",
		"status":    status,
		"kubectl":   output,
	})
}
