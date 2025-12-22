package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"dbaas-orcastrator/internal/models"
)

// HealthCheck - Health check endpoint
func (h *KubeDBHandler) MonitoringHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:    "success",
		Message:   "DBaaS Monitoring API is running",
		Timestamp: time.Now().Unix(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}