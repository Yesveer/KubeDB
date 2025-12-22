package installer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type k8sItem struct {
	Spec map[string]interface{} `json:"spec"`
}

type k8sList struct {
	Items []k8sItem `json:"items"`
}

func GetDBVersions(kubeconfigPath, dbType string) ([]string, error) {

	resource := map[string]string{
		"mongo":      "mongodbversions",
		"mysql":      "mysqlversions",
		"postgres":   "postgresversions",
		"redis":      "redisversions",
		"clickhouse": "clickhouseversions",
		"kafka":      "kfversion",
	}

	res, ok := resource[dbType]
	if !ok {
		return nil, fmt.Errorf("unsupported dbType: %s", dbType)
	}

	cmd := exec.Command(
		"kubectl",
		"get", res,
		"-o", "json",
	)

	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+kubeconfigPath,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var list k8sList
	if err := json.Unmarshal(out.Bytes(), &list); err != nil {
		return nil, err
	}

	versionSet := map[string]bool{}
	var versions []string

	for _, item := range list.Items {

		version, _ := item.Spec["version"].(string)
		if version == "" {
			continue
		}

		// 🔥 Mongo: only Official
		if dbType == "mongo" {
			if dist, ok := item.Spec["distribution"].(string); ok {
				if strings.ToLower(dist) != "official" {
					continue
				}
			}
		}

		// 🔥 Redis: skip valkey
		if dbType == "redis" {
			if strings.Contains(strings.ToLower(version), "valkey") {
				continue
			}
		}

		// ✅ Kafka & ClickHouse → NO FILTER (like redis but simpler)

		if !versionSet[version] {
			versionSet[version] = true
			versions = append(versions, version)
		}
	}

	return versions, nil
}
