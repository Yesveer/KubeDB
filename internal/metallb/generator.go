package metallb

import (
	"os"
	"text/template"
)

const metallbTemplate = `
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: {{ .PoolName }}
  namespace: metallb-system
spec:
  addresses:
{{- range .Addresses }}
  - {{ . }}
{{- end }}
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: {{ .AdvertisementName }}
  namespace: metallb-system
spec:
  ipAddressPools:
    - {{ .PoolName }}
`

type TemplateData struct {
	PoolName          string
	Addresses         []string
	AdvertisementName string
}

func GenerateYAML(path string, data TemplateData) error {

	tmpl, err := template.New("metallb").Parse(metallbTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
