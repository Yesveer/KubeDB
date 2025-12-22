// package handlers

// import (
// 	"encoding/json"
// 	"net/http"

// 	"dbaas-orcastrator/internal/repository"
// )

// func (h *KubeDBHandler) ListDBaaS(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	data, err := repository.ListAllDBaaS()
// 	if err != nil {
// 		http.Error(w, "DB query failed: "+err.Error(), 500)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"count": len(data),
// 		"data":  data,
// 	})
// }

package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/models"
	"dbaas-orcastrator/internal/repository"
)

// func (h *KubeDBHandler) ListDBaaS(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// 🔹 Read query params
// 	domain := r.URL.Query().Get("domain")
// 	project := r.URL.Query().Get("project")

// 	if domain == "" || project == "" {
// 		http.Error(w, "domain and project are required", http.StatusBadRequest)
// 		return
// 	}

// 	// 🔹 Call filtered repo
// 	data, err := repository.ListDBaaSByDomainProject(domain, project)
// 	if err != nil {
// 		http.Error(w, "DB query failed: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"count": len(data),
// 		"data":  data,
// 	})
// }

func (h *KubeDBHandler) ListDBaaS(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	domain := r.URL.Query().Get("domain")
	project := r.URL.Query().Get("project")

	var (
		data []models.DBaaSRecord
		err  error
	)

	// 🔹 ALL DBaaS
	if domain == "" && project == "" {
		data, err = repository.ListAllDBaaS()
	} else if domain != "" && project != "" {
		data, err = repository.ListDBaaSByDomainProject(domain, project)
	} else {
		http.Error(w, "both domain and project are required", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "DB query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(data),
		"data":  data,
	})
}
