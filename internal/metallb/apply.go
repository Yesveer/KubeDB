package metallb

import (
	"fmt"
	"os"
	"os/exec"
)

func ApplyYAML(kubeconfigPath, yamlPath string) error {

	cmd := exec.Command("kubectl", "apply", "-f", yamlPath)

	// 🔥 IMPORTANT: force kubeconfig
	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+kubeconfigPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubectl apply failed: %s", out)
	}

	return nil
}
