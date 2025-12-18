package backup

import (
	"fmt"
	"log"
	"os/exec"
)

func InstallVeleroViaScript(s3Url, accessKey, secretKey string) error {

	cmd := exec.Command(
		"bash",
		"./scripts/velero_install.sh",
		s3Url,
		accessKey,
		secretKey,
	)

	out, err := cmd.CombinedOutput()
	log.Println("VELERO SCRIPT OUTPUT:\n", string(out))

	if err != nil {
		return fmt.Errorf("velero install failed:\n%s", out)
	}

	return nil
}
