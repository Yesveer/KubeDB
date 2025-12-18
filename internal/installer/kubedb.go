package installer

import (
	"os"
	"os/exec"
)

func InstallKubeDB() error {
	scriptPath := os.Getenv("INSTALL_SCRIPT_PATH")
	if scriptPath == "" {
		scriptPath = "./scripts/install_kubedb.sh"
	}

	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
