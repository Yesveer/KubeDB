package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/services"
)

// GetPrometheusURL - Download kubeconfig & return Prometheus endpoint
func (h *KubeDBHandler) GetPrometheusURL(w http.ResponseWriter, r *http.Request) {

	// 1️⃣ Read query params
	domain := r.URL.Query().Get("domain")
	project := r.URL.Query().Get("project")
	cluster := r.URL.Query().Get("cluster")

	if domain == "" || project == "" || cluster == "" {
		http.Error(w, "domain, project and cluster are required", http.StatusBadRequest)
		return
	}

	// 2️⃣ Get Bearer token
	token, err := getBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 3️⃣ Download kubeconfig
	kubeconfigPath := fmt.Sprintf(
		"/tmp/kubeconfig-%s-%s-%s.yaml",
		domain, project, cluster,
	)

	if err := kubeconfig.Download(
		h.Cfg.CompassBaseURL,
		token,
		domain,
		project,
		cluster,
		kubeconfigPath,
	); err != nil {
		http.Error(w, "kubeconfig download failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4️⃣ Discover Prometheus URL via kubectl
	kubectlSvc := services.NewKubectlService(kubeconfigPath)
	promURL, err := kubectlSvc.DiscoverPrometheusURL()
	if err != nil {
		http.Error(w, "prometheus discovery failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5️⃣ Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":        "success",
		"prometheusURL": promURL,
	})
}
