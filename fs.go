package docs_blog_engine

import (
	"os"
)

func ensureFoldersExist(folders []string) error {
	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
