package installer

import (
	"bytes"
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

	cmd.Env = append(os.Environ(), "KUBECONFIG="+os.Getenv("KUBECONFIG_PATH"))

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
