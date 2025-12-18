package kubeconfig

import (
	"encoding/base64"
	"os"
)

func EncodeFileToBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
