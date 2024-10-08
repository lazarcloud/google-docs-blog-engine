package backup

import (
	"os"
	"path/filepath"
	"time"

	files "github.com/lazarcloud/google-docs-blog-engine/fs"
	"github.com/lazarcloud/google-docs-blog-engine/globals"
)

func CreateBackup() error {
	timestamp := time.Now().Format("2006-01-02-15-04-05")

	path := filepath.Join(globals.BackupDir, timestamp)

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = files.CopyDir(globals.StaticDir, path)
	return err
}
