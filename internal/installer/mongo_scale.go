package installer

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func ScaleMongo(name, namespace string, replicas int, storage string) error {

	kubeconfig := os.Getenv("KUBECONFIG_PATH")

	log.Println("🔄 Scaling MongoDB")
	log.Println("➡️ Name:", name)
	log.Println("➡️ Namespace:", namespace)
	log.Println("➡️ Replicas:", replicas)
	log.Println("➡️ Storage:", storage)

	patch := fmt.Sprintf(`{
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
	}`, replicas, storage)

	log.Println("📄 Patch payload:", patch)

	cmd := exec.Command(
		"kubectl", "patch", "mongodb", name,
		"-n", namespace,
		"--type=merge",
		"-p", patch,
	)

	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+kubeconfig,
	)

	out, err := cmd.CombinedOutput()
	log.Println("📤 kubectl output:\n", string(out))

	if err != nil {
		log.Println("❌ ScaleMongo failed:", err)
		return fmt.Errorf("kubectl patch failed: %w", err)
	}

	log.Println("✅ MongoDB scaled successfully")
	return nil
}
