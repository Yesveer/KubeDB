package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/repository"
)

func (h *KubeDBHandler) ListMetalLBPools(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}

	pools, err := repository.ListMetalLBPools()
	if err != nil {
		http.Error(w, "DB query failed", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(pools),
		"data":  pools,
	})
}
