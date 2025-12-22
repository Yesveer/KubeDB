package services

import (
	"fmt"
	"os/exec"
	"strings"
)

type KubectlService struct {
	KubeconfigPath string
}

// NewKubectlService - Initialize kubectl service
func NewKubectlService(kubeconfigPath string) *KubectlService {
	return &KubectlService{
		KubeconfigPath: kubeconfigPath,
	}
}

// runCommand - Execute kubectl command with kubeconfig
func (ks *KubectlService) runCommand(args ...string) (string, error) {
	// Add kubeconfig flag
	cmdArgs := append([]string{"--kubeconfig", ks.KubeconfigPath}, args...)
	
	// Execute command
	cmd := exec.Command("kubectl", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("kubectl failed: %s, output: %s", err.Error(), string(output))
	}
	
	return strings.TrimSpace(string(output)), nil
}

// PatchPrometheusToNodePort - Convert Prometheus service to NodePort
func (ks *KubectlService) PatchPrometheusToNodePort() error {
	_, err := ks.runCommand(
		"patch", "svc", "kube-prom-stack-kube-prome-prometheus",
		"-n", "monitoring",
		"-p", `{"spec":{"type":"NodePort"}}`,
	)
	return err
}

// GetNodeIP - Get first node's internal IP
func (ks *KubectlService) GetNodeIP() (string, error) {
	output, err := ks.runCommand(
		"get", "nodes",
		"-o", `jsonpath={.items[0].status.addresses[?(@.type=="InternalIP")].address}`,
	)
	if err != nil {
		return "", err
	}
	
	if output == "" {
		return "", fmt.Errorf("no node IP found")
	}
	
	return output, nil
}

// GetPrometheusNodePort - Get Prometheus service NodePort
func (ks *KubectlService) GetPrometheusNodePort() (string, error) {
	output, err := ks.runCommand(
		"get", "svc", "kube-prom-stack-kube-prome-prometheus",
		"-n", "monitoring",
		"-o", `jsonpath={.spec.ports[?(@.port==9090)].nodePort}`,
	)
	if err != nil {
		return "", err
	}
	
	if output == "" {
		return "", fmt.Errorf("no NodePort found for Prometheus")
	}
	
	return output, nil
}

// DiscoverPrometheusURL - Complete flow to discover Prometheus endpoint
func (ks *KubectlService) DiscoverPrometheusURL() (string, error) {
	
	// Step 1: Patch service to NodePort (ignore error if already NodePort)
	_ = ks.PatchPrometheusToNodePort()
	
	// Step 2: Get Node IP
	nodeIP, err := ks.GetNodeIP()
	if err != nil {
		return "", fmt.Errorf("failed to get node IP: %w", err)
	}
	
	// Step 3: Get NodePort
	nodePort, err := ks.GetPrometheusNodePort()
	if err != nil {
		return "", fmt.Errorf("failed to get Prometheus NodePort: %w", err)
	}
	
	// Step 4: Build URL
	prometheusURL := fmt.Sprintf("http://%s:%s", nodeIP, nodePort)
	
	return prometheusURL, nil
}