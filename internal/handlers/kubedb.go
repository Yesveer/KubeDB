package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"dbaas-orcastrator/internal/config"
	"dbaas-orcastrator/internal/installer"
	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/repository"
)

type KubeDBHandler struct {
	Cfg *config.Config
}

type InstallRequest struct {
	Domain  string `json:"domain"`
	Project string `json:"project"`
	Cluster string `json:"cluster"`
}

// Helper: Bearer token extract
func getBearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("Authorization header missing")
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("Authorization header must be Bearer token")
	}

	return strings.TrimPrefix(auth, "Bearer "), nil
}

// POST /install/kubedb
func (h *KubeDBHandler) Install(w http.ResponseWriter, r *http.Request) {

	// 1️⃣ Parse request
	var req InstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 2️⃣ Get Bearer token
	token, err := getBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 3️⃣ Download kubeconfig
	if err := kubeconfig.Download(
		h.Cfg.CompassBaseURL,
		token,
		req.Domain,
		req.Project,
		req.Cluster,
		h.Cfg.KubeconfigPath,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4️⃣ Encode kubeconfig
	kubeconfigB64, err := kubeconfig.EncodeFileToBase64(h.Cfg.KubeconfigPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 🔥 5️⃣ SAVE STATUS = INSTALLING (IMPORTANT)
	if err := repository.SaveDBaaSRecord(
		req.Domain,
		req.Project,
		req.Cluster,
		"INSTALLING",
		kubeconfigB64,
	); err != nil {
		http.Error(w, "Mongo save failed: "+err.Error(), 500)
		return
	}

	// 6️⃣ Run install
	if err := installer.InstallKubeDB(); err != nil {

		// ❌ INSTALL FAILED → update status
		_ = repository.SaveDBaaSRecord(
			req.Domain,
			req.Project,
			req.Cluster,
			"FAILED",
			kubeconfigB64,
		)

		http.Error(w, "KubeDB install failed: "+err.Error(), 500)
		return
	}

	// ✅ 7️⃣ INSTALL SUCCESS → update status
	if err := repository.SaveDBaaSRecord(
		req.Domain,
		req.Project,
		req.Cluster,
		"INSTALLED",
		kubeconfigB64,
	); err != nil {
		http.Error(w, "Mongo save failed: "+err.Error(), 500)
		return
	}

	// 8️⃣ Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "KubeDB installed successfully",
	})
}
