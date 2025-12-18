package licence

import (
	"encoding/base64"
)

func EncodeToBase64(license string) string {
	return base64.StdEncoding.EncodeToString([]byte(license))
}
