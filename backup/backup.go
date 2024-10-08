package backup

func CreateBackup() error {
	timestamp := time.Now().Format("2006-01-02-15-04-05")

	path := filepath.Join(globals.BackupDir, timestamp)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err := files.CopyDir(globals.StaticDir, path)
	if err != nil {
		return err
	}
}