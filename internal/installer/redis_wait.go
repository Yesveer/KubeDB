package installer

import (
	"os"
	"os/exec"
	"time"
)

func WaitForRedisReady(kubeconfig, name, namespace string) error {
	for {
		cmd := exec.Command(
			"kubectl",
			"get", "redis.kubedb.com", name,
			"-n", namespace,
			"-o=jsonpath={.status.phase}",
		)
		cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfig)

		out, _ := cmd.Output()
		if string(out) == "Ready" {
			return nil
		}

		time.Sleep(10 * time.Second)
	}
}
