package installer

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"time"

	"dbaas-orcastrator/internal/models"
	"dbaas-orcastrator/internal/repository"
)

func CreateVeleroBackup(
	domain, project, cluster string,
	dbName string,
	namespace string,
	backupName string,
) {

	log.Println("📦 Starting Velero backup:", backupName)

	cmd := exec.Command(
		"velero",
		"backup", "create", backupName,
		"--include-namespaces", namespace,
		"--default-volumes-to-fs-backup",
		"--wait",
	)

	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+os.Getenv("KUBECONFIG_PATH"),
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Println("❌ Velero backup failed:", err)
		return
	}

	// 🔥 BACKUP DESCRIBE (JSON)
	descCmd := exec.Command(
		"velero", "backup", "get", backupName, "-o", "json",
	)
	descCmd.Env = cmd.Env

	var desc bytes.Buffer
	descCmd.Stdout = &desc
	if err := descCmd.Run(); err != nil {
		log.Println("❌ Failed to describe backup:", err)
		return
	}

	var v struct {
		Status struct {
			Phase               string    `json:"phase"`
			StartTimestamp      time.Time `json:"startTimestamp"`
			CompletionTimestamp time.Time `json:"completionTimestamp"`
		} `json:"status"`
		Spec struct {
			TTL string `json:"ttl"`
		} `json:"spec"`
	}

	if err := json.Unmarshal(desc.Bytes(), &v); err != nil {
		log.Println("❌ JSON parse failed:", err)
		return
	}

	backup := models.BackupInfo{
		BackupName:  backupName,
		Status:      v.Status.Phase,
		StartedAt:   v.Status.StartTimestamp,
		CompletedAt: v.Status.CompletionTimestamp,
		TTL:         v.Spec.TTL,
	}

	// 🔥 DB UPDATE
	if err := repository.AppendMongoBackup(
		domain,
		project,
		cluster,
		dbName, // 🔥 mongo-rs
		backup,
	); err != nil {
		log.Println("❌ Failed to store backup in DB:", err)
		return
	}

	log.Println("✅ Backup stored in DB:", backupName)
}
