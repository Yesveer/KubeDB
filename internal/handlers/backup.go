package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/backup"
	"dbaas-orcastrator/internal/repository"
)

type BackupRequest struct {
	Domain  string `json:"domain"`
	Project string `json:"project"`
	Cluster string `json:"cluster"`

	S3Url        string `json:"s3Url"`
	AWSAccessKey string `json:"aws_access_key_id"`
	AWSSecretKey string `json:"aws_secret_access_key"`
}

func (h *KubeDBHandler) EnableBackup(w http.ResponseWriter, r *http.Request) {

	// 1️⃣ Parse request
	var req BackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 2️⃣ Basic validation
	if req.S3Url == "" || req.AWSAccessKey == "" || req.AWSSecretKey == "" {
		http.Error(w, "s3Url / aws credentials missing", http.StatusBadRequest)
		return
	}

	// 3️⃣ Install Velero via SHELL SCRIPT ✅
	if err := backup.InstallVeleroViaScript(
		req.S3Url,
		req.AWSAccessKey,
		req.AWSSecretKey,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4️⃣ Update DB
	if err := repository.EnableBackup(
		req.Domain,
		req.Project,
		req.Cluster,
		req.S3Url,
		req.AWSAccessKey,
		req.AWSSecretKey,
	); err != nil {
		http.Error(w, "Mongo update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5️⃣ Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Backup enabled successfully",
	})
}
