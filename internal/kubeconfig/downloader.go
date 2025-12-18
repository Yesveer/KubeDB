package kubeconfig

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(
	baseURL, token, domain, project, cluster, path string,
) error {

	url := fmt.Sprintf(
		"%s/api/compass/v1/domain/%s/project/%s/cluster/%s/kubeconfig",
		baseURL, domain, project, cluster,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("organisation-name", domain)
	req.Header.Set("external-project", project)
	req.Header.Set("project-name", project)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("kubeconfig download failed: %d", resp.StatusCode)
	}

	out, err := os.Create(path) // overwrite
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
