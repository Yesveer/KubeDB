package installer

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GetClusterUID() (string, error) {

	cmd := exec.Command(
		"kubectl",
		"get", "ns", "kube-system",
		"-o=jsonpath={.metadata.uid}",
	)

	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+os.Getenv("KUBECONFIG_PATH"),
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	clusterUID := strings.TrimSpace(out.String())

	// 🔥 LOG CLUSTER UID
	log.Println("✅ Kubernetes Cluster UID:", clusterUID)

	return clusterUID, nil
}
