package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetFirstNodeIP(kubeconfig string) (string, error) {
	cmd := exec.Command(
		"kubectl",
		"get", "nodes",
		"-o=jsonpath={.items[0].status.addresses[?(@.type==\"InternalIP\")].address}",
	)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfig)

	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// func AddIPToNode(nodeIP, lbIP, user, password string) error {
// 	cmd := exec.Command(
// 		"sshpass", "-p", password,
// 		"ssh",
// 		"-o", "StrictHostKeyChecking=no",
// 		fmt.Sprintf("%s@%s", user, nodeIP),
// 		fmt.Sprintf("sudo ip addr add %s/32 dev ens3 || true", lbIP),
// 	)

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()
// }

func AddIPToNode(nodeIP, lbIP string) error {

	user := os.Getenv("NODE_USER")
	pass := os.Getenv("NODE_PASSWORD")

	if user == "" || pass == "" {
		return fmt.Errorf("NODE_USER or NODE_PASSWORD not set")
	}

	cmdStr := fmt.Sprintf(
		`sshpass -p '%s' ssh -o StrictHostKeyChecking=no %s@%s "echo '%s' | sudo -S ip addr add %s/32 dev ens3"`,
		pass,
		user,
		nodeIP,
		pass,
		lbIP,
	)

	cmd := exec.Command("bash", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ssh failed: %s", string(out))
	}

	return nil
}
