package metallb

import (
	"fmt"
	"os"
	"os/exec"
)

func DeletePool(kubeconfigPath, poolName string) error {

	cmd := exec.Command(
		"kubectl",
		"delete",
		"ipaddresspool",
		poolName,
		"-n", "metallb-system",
	)

	// 🔥 Force kubeconfig
	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+kubeconfigPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubectl delete failed: %s", out)
	}

	return nil
}
