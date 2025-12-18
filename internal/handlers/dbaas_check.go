package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/repository"
)

func (h *KubeDBHandler) CheckDBaaS(w http.ResponseWriter, r *http.Request) {

	domain := r.URL.Query().Get("domain")
	project := r.URL.Query().Get("project")
	cluster := r.URL.Query().Get("cluster")

	if domain == "" || project == "" || cluster == "" {
		http.Error(w, "domain, project, cluster are required", http.StatusBadRequest)
		return
	}

	record, err := repository.FindDBaaSRecord(domain, project, cluster)
	if err != nil {
		http.Error(w, "Mongo query failed: "+err.Error(), 500)
		return
	}

	if record == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exists":  false,
			"message": "DBaaS record not found",
		})
		return
	}

	// Found
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exists": true,
		"data": map[string]interface{}{
			"domain":         record.Domain,
			"project":        record.Project,
			"cluster":        record.Cluster,
			"status":         record.Status,
			"backup_enabled": record.BackupEnabled,
			"updatedAt":      record.UpdatedAt,
		},
	})
}
