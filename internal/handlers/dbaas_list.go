package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/repository"
)

func (h *KubeDBHandler) ListDBaaS(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := repository.ListAllDBaaS()
	if err != nil {
		http.Error(w, "DB query failed: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(data),
		"data":  data,
	})
}
