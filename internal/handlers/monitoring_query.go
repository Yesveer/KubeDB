package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/models"
	"dbaas-orcastrator/internal/services"
)

// Helper: Extract domain, project, cluster from request
func getClusterInfo(r *http.Request) (string, string, string, error) {
	domain := r.URL.Query().Get("domain")
	project := r.URL.Query().Get("project")
	cluster := r.URL.Query().Get("cluster")
	
	if domain == "" || project == "" || cluster == "" {
		return "", "", "", fmt.Errorf("domain, project, and cluster query parameters are required")
	}
	
	return domain, project, cluster, nil
}

// QueryMetrics - Execute custom Prometheus query
func (h *KubeDBHandler) QueryMetrics(w http.ResponseWriter, r *http.Request) {
	
	// 1. Get query parameter
	query := r.URL.Query().Get("query")
	if query == "" {
		response := models.MetricsResponse{
			Status: "error",
			Error:  "query parameter is required",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 2. Get cluster info
	domain, project, cluster, err := getClusterInfo(r)
	if err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 3. Get Bearer token
	token, err := getBearerToken(r)
	if err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 4. Download kubeconfig
	kubeconfigPath := fmt.Sprintf("/tmp/kubeconfig-%s-%s-%s.yaml", domain, project, cluster)
	if err := kubeconfig.Download(
		h.Cfg.CompassBaseURL,
		token,
		domain,
		project,
		cluster,
		kubeconfigPath,
	); err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  fmt.Sprintf("kubeconfig download failed: %s", err.Error()),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 5. Discover Prometheus URL using kubectl
	kubectlSvc := services.NewKubectlService(kubeconfigPath)
	prometheusURL, err := kubectlSvc.DiscoverPrometheusURL()
	if err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  fmt.Sprintf("Prometheus discovery failed: %s", err.Error()),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 6. Initialize Prometheus service with discovered URL
	promService := services.NewPrometheusService(prometheusURL)
	
	// 7. Execute query
	result, err := promService.Query(r.Context(), query)
	if err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 8. Success response
	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// QueryRangeMetrics - Execute Prometheus range query
func (h *KubeDBHandler) QueryRangeMetrics(w http.ResponseWriter, r *http.Request) {
	
	// 1. Get query parameter
	query := r.URL.Query().Get("query")
	if query == "" {
		response := models.MetricsResponse{
			Status: "error",
			Error:  "query parameter is required",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 2. Get cluster info
	domain, project, cluster, err := getClusterInfo(r)
	if err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 3. Get Bearer token
	token, err := getBearerToken(r)
	if err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 4. Download kubeconfig
	kubeconfigPath := fmt.Sprintf("/tmp/kubeconfig-%s-%s-%s.yaml", domain, project, cluster)
	if err := kubeconfig.Download(
		h.Cfg.CompassBaseURL,
		token,
		domain,
		project,
		cluster,
		kubeconfigPath,
	); err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  fmt.Sprintf("kubeconfig download failed: %s", err.Error()),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 5. Discover Prometheus URL using kubectl
	kubectlSvc := services.NewKubectlService(kubeconfigPath)
	prometheusURL, err := kubectlSvc.DiscoverPrometheusURL()
	if err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  fmt.Sprintf("Prometheus discovery failed: %s", err.Error()),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 6. Parse time parameters
	end := time.Now()
	start := end.Add(-1 * time.Hour)
	step := 30 * time.Second
	
	if startParam := r.URL.Query().Get("start"); startParam != "" {
		if t, err := time.Parse(time.RFC3339, startParam); err == nil {
			start = t
		}
	}
	if endParam := r.URL.Query().Get("end"); endParam != "" {
		if t, err := time.Parse(time.RFC3339, endParam); err == nil {
			end = t
		}
	}
	
	// 7. Initialize Prometheus service with discovered URL
	promService := services.NewPrometheusService(prometheusURL)
	
	// 8. Execute query
	result, err := promService.QueryRange(r.Context(), query, start, end, step)
	if err != nil {
		response := models.MetricsResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 9. Success response
	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}