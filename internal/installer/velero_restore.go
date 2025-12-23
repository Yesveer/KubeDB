package installer

import (
	"log"
	"os"
	"os/exec"
)

func CreateVeleroRestore(restoreName, backupName string) error {

	log.Println("♻️ Starting Velero restore:", restoreName)

	cmd := exec.Command(
		"velero",
		"restore", "create", restoreName,
		"--from-backup", backupName,
		"--wait",
	)

	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+os.Getenv("KUBECONFIG_PATH"),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("❌ Velero restore failed:", string(out))
		return err
	}

	log.Println("✅ Velero restore completed:", restoreName)
	return nil
}
