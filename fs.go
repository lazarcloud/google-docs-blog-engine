package docs_blog_engine

import (
	"io"
	"os"
	"path/filepath"
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

func clearDirectories(folders []string) error {
	for _, folder := range folders {
		if err := os.RemoveAll(folder); err != nil {
			return err
		}
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the file contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Preserve file permissions
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

// copyDir recursively copies a directory tree, attempting to preserve permissions.
func copyDir(src, dst string) error {
	// Get properties of source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// Read the directory contents
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Loop through directory contents
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// If it's a directory, recursively copy
		if entry.IsDir() {
			err := copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Copy file
			err := copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
