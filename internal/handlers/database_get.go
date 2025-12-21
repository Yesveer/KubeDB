package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/repository"
)

func (h *KubeDBHandler) GetDatabases(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	domain := r.URL.Query().Get("domain")
	project := r.URL.Query().Get("project")
	name := r.URL.Query().Get("name")

	if domain == "" || project == "" {
		http.Error(w, "domain and project are required", http.StatusBadRequest)
		return
	}

	// 🔹 Single DB
	if name != "" {
		db, err := repository.GetDatabaseByName(domain, project, name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if db == nil {
			json.NewEncoder(w).Encode(map[string]any{
				"exists": false,
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"exists": true,
			"data":   db,
		})
		return
	}

	// 🔹 All DBs
	list, err := repository.GetDatabases(domain, project)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"count": len(list),
		"data":  list,
	})
}
