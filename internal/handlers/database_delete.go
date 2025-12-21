package handlers

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/repository"
)

type DeleteDatabaseRequest struct {
	Domain    string `json:"domain"`
	Project   string `json:"project"`
	Name      string `json:"name"`
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`
}

func (h *KubeDBHandler) DeleteDatabase(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	token, err := getBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 1️⃣ Download kubeconfig
	if err := kubeconfig.Download(
		h.Cfg.CompassBaseURL,
		token,
		req.Domain,
		req.Project,
		req.Cluster,
		h.Cfg.KubeconfigPath,
	); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 1️⃣ Delete namespace
	cmd := exec.Command("kubectl", "delete", "ns", req.Namespace)
	cmd.Env = append(cmd.Env, "KUBECONFIG="+h.Cfg.KubeconfigPath)
	_, _ = cmd.CombinedOutput() // ignore error if already deleted

	// 2️⃣ Delete DB entry
	if err := repository.DeleteDatabase(req.Domain, req.Project, req.Name); err != nil {
		http.Error(w, "DB delete failed: "+err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Database deleted successfully",
	})
}
