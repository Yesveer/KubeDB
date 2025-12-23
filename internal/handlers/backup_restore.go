package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"dbaas-orcastrator/internal/installer"
	"dbaas-orcastrator/internal/kubeconfig"
	"dbaas-orcastrator/internal/repository"
)

type RestoreBackupRequest struct {
	Domain        string `json:"domain"`
	Project       string `json:"project"`
	Cluster       string `json:"cluster"`       // SOURCE CLUSTER
	Name          string `json:"name"`          // DB NAME
	TargetCluster string `json:"targetCluster"` // TARGET CLUSTER
}

func (h *KubeDBHandler) RestoreBackup(w http.ResponseWriter, r *http.Request) {

	var req RestoreBackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	token, err := getBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// 🔹 1. SOURCE CLUSTER kubeconfig
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

	// 🔹 2. GET LATEST COMPLETED BACKUP FROM DB
	backup, err := repository.GetLatestCompletedBackup(
		req.Domain,
		req.Project,
		req.Cluster,
		req.Name,
	)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	// 🔹 3. TARGET CLUSTER kubeconfig
	if err := kubeconfig.Download(
		h.Cfg.CompassBaseURL,
		token,
		req.Domain,
		req.Project,
		req.TargetCluster,
		h.Cfg.KubeconfigPath,
	); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	restoreName := req.Name + "-restore-" + time.Now().Format("20060102-150405")

	// 🔥 4. ASYNC RESTORE
	go installer.CreateVeleroRestore(
		restoreName,
		backup.BackupName,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":        "Restore started",
		"backupName":    backup.BackupName,
		"restoreName":   restoreName,
		"sourceCluster": req.Cluster,
		"targetCluster": req.TargetCluster,
	})
}
