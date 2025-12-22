package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/installer"
	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/repository"
)

type UpgradeDBRequest struct {
	Domain        string `json:"domain"`
	Project       string `json:"project"`
	Cluster       string `json:"cluster"`
	DBType        string `json:"dbType"`
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	TargetVersion string `json:"targetVersion"`
}

func (h *KubeDBHandler) UpgradeDatabase(w http.ResponseWriter, r *http.Request) {

	var req UpgradeDBRequest
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

	// 🔥 Run upgrade async
	go func() {

		err := installer.UpgradeDatabase(
			req.DBType,
			req.Name,
			req.Namespace,
			req.TargetVersion,
		)
		if err != nil {
			return
		}

		// ✅ Update DB version after successful apply
		_ = repository.UpdateDBVersion(
			req.Domain,
			req.Project,
			req.Cluster,
			req.Name,
			req.TargetVersion,
		)
	}()

	json.NewEncoder(w).Encode(map[string]string{
		"status":        "Upgrade started",
		"targetVersion": req.TargetVersion,
	})
}
