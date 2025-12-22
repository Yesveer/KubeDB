package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func UpgradeDatabase(dbType, name, namespace, targetVersion string) error {

	kindMap := map[string]string{
		"mongo":      "MongoDBOpsRequest",
		"mysql":      "MySQLOpsRequest",
		"postgres":   "PostgresOpsRequest",
		"redis":      "RedisOpsRequest",
		"clickhouse": "ClickHouseOpsRequest",
		"kafka":      "KafkaOpsRequest",
	}

	kind, ok := kindMap[dbType]
	if !ok {
		return fmt.Errorf("unsupported dbType: %s", dbType)
	}

	yaml := fmt.Sprintf(`
apiVersion: ops.kubedb.com/v1alpha1
kind: %s
metadata:
  name: %s-upgrade-%s
  namespace: %s
spec:
  type: UpdateVersion
  databaseRef:
    name: %s
  updateVersion:
    targetVersion: "%s"
  timeout: 30m
  apply: IfReady
`, kind, name, targetVersion, namespace, name, targetVersion)

	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+os.Getenv("KUBECONFIG_PATH"),
	)
	cmd.Stdin = strings.NewReader(yaml)

	_, err := cmd.CombinedOutput()
	return err
}
