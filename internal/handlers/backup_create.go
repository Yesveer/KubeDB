package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"dbaas-orcastrator/internal/installer"
	"dbaas-orcastrator/internal/kubeconfig"
)

type CreateBackupRequest struct {
	Domain    string `json:"domain"`
	Project   string `json:"project"`
	Cluster   string `json:"cluster"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (h *KubeDBHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {

	var req CreateBackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	token, err := getBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 1️⃣ kubeconfig download
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

	backupName := req.Name + "-" + time.Now().Format("20060102-150405")

	// 🔥 ASYNC BACKUP
	go installer.CreateVeleroBackup(
		req.Domain,
		req.Project,
		req.Cluster,
		req.Name, // ✅ DB NAME (mongo-rs)
		req.Namespace,
		backupName, // ✅ BACKUP NAME
	)

	json.NewEncoder(w).Encode(map[string]string{
		"status":      "Backup started",
		"backup_name": backupName,
	})
}
