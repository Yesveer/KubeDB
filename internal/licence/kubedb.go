package licence

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func Generate(name, email, clusterUID string) error {

	data := url.Values{}
	data.Set("name", name)
	data.Set("email", email)
	data.Set("product", "kubedb-enterprise")
	data.Set("cluster", clusterUID)
	data.Set("coupon", "")
	data.Set("tos", "true")

	req, err := http.NewRequest(
		"POST",
		"https://appscode.com/issue-license",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("license API failed: %d", resp.StatusCode)
	}

	return nil
}
