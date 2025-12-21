package installer

import (
	"fmt"
	"os"
	"os/exec"
)

func ScaleMongo(name, namespace string, replicas int, storage string) error {

	kubeconfig := os.Getenv("KUBECONFIG_PATH")

	cmd := exec.Command(
		"kubectl", "patch", "mongodb", name,
		"-n", namespace,
		"--type=merge",
		"-p",
		fmt.Sprintf(`{
		  "spec": {
		    "replicas": %d,
		    "storage": {
		      "resources": {
		        "requests": {
		          "storage": "%s"
		        }
		      }
		    }
		  }
		}`, replicas, storage),
	)

	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfig)

	_, err := cmd.CombinedOutput()
	return err
}
