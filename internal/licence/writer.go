package licence

import (
	"os"
)

func WriteToFile(path, license string) error {
	return os.WriteFile(path, []byte(license), 0644)
}
