package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/repository"
)

func (h *KubeDBHandler) ListMetalLBPools(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ✅ Query params
	domain := r.URL.Query().Get("domain")
	project := r.URL.Query().Get("project")
	cluster := r.URL.Query().Get("cluster")

	pools, err := repository.ListMetalLBPoolsFiltered(
		domain,
		project,
		cluster,
	)
	if err != nil {
		http.Error(w, "DB query failed: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(pools),
		"data":  pools,
	})
}
