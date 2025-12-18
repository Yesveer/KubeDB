package installer

import (
	"bytes"
	"os"
	"os/exec"
)

func IsKubeDBInstalled() (bool, string, error) {

	cmd := exec.Command(
		"kubectl",
		"get", "pods",
		"--all-namespaces",
		"-l", "app.kubernetes.io/instance=kubedb",
	)

	cmd.Env = append(os.Environ(), "KUBECONFIG="+os.Getenv("KUBECONFIG_PATH"))

	var out bytes.Buffer
	var errBuf bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		return false, errBuf.String(), err
	}

	// If output has more than header line, pods exist
	lines := bytes.Count(out.Bytes(), []byte("\n"))
	if lines > 1 {
		return true, out.String(), nil
	}

	return false, out.String(), nil
}
