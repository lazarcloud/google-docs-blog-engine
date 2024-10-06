package docs_blog_engine

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Post struct {
	Title       string `json:"title"`
	ID          string `json:"id"`
	Date        string `json:"date"`
	ContentHTML string `json:"content_html"`
	ContentMD   string `json:"content_md"`
}

var posts []Post

func getDoc(docID string, format string) (string, error) {
	// download google docs file as html or md
	// use the srv variable to access the drive service
	// use the docID to get the file

	// determine the export MIME type based on the format
	var mimeType string
	switch format {
	case "html":
		mimeType = "text/html"
	case "md":
		mimeType = "text/markdown"
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	// get the file
	file, err := srv.Files.Export(docID, mimeType).Download()
	if err != nil {
		return "", err
	}
	defer file.Body.Close()

	// read the file
	buf := make([]byte, 1024)
	content := ""
	for {
		n, err := file.Body.Read(buf)
		content += string(buf[:n])
		if err != nil {
			break
		}
	}
	return content, nil
}

func getPosts() error {
	dates := []string{}
	// print all fields of the folder

	// Get the last modified time of the folder
	folder, err := srv.Files.Get(string(folderID)).Do()
	if err != nil {
		return err
	}
	fmt.Printf("Folder last modified time: %s\n", folder.ModifiedDate)

	dates = append(dates, folder.ModifiedDate)

	fileList, err := srv.Files.List().Q(fmt.Sprintf("'%s' in parents", folderID)).Do()
	if err != nil {
		// log.Fatalf("Unable to retrieve files in folder: %v", err.Error())
		return err
	}

	for _, file := range fileList.Items {
		dates = append(dates, file.ModifiedDate)
	}

	newestModified := dates[0]
	for _, date := range dates {
		if date > newestModified {
			newestModified = date
		}
	}

	if newestModified == lastChanged {
		fmt.Println("No changes in the folder")
		return nil
	}

	fmt.Println("Changes detected in the folder")

	newPosts := []Post{}

	// Clear the markdown folder with images
	path := "./app/src/content/blog"
	err = os.RemoveAll(path)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	path = "./app/public/images"
	err = os.RemoveAll(path)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	for _, file := range fileList.Items {
		fmt.Printf("Name: %s, ID: %s\n", file.Title, file.Id)
		fmt.Println(file.CreatedDate)
		fmt.Println(file.MimeType)
		html, err := getDoc(file.Id, "html")
		if err != nil {
			return err
		}
		md, err := getDoc(file.Id, "md")
		if err != nil {
			return err
		}
		newPosts = append(newPosts, Post{
			Title:       file.Title,
			ID:          file.Id,
			Date:        file.CreatedDate,
			ContentHTML: html,
			ContentMD:   md,
		})

		fmt.Println(md)

		firstLine := strings.Split(md, "\n")[0]
		if !strings.HasPrefix(firstLine, "DESCRIPTION ") {
			firstLine = "No description"
		} else {
			firstLine = strings.TrimPrefix(firstLine, "DESCRIPTION ")
			md = strings.Join(strings.Split(md, "\n")[1:], "\n")
		}

		fileID := strings.ToLower(strings.ReplaceAll(file.Title, " ", "-"))

		// Remove the first image from the markdown content
		md = removeFirstImage(md)
		// "2024-10-05T15:41:20.896Z" -> 10 Jan 2024
		parsedDate, err := time.Parse(time.RFC3339, file.CreatedDate)
		if err != nil {
			return err
		}

		// Format the parsed date
		formattedDate := parsedDate.Format("02 Jan 2006")

		toAppend := fmt.Sprintf(`---
title: '%s'
description: '%s'
pubDate: '%s'
heroImage: '/images/%s-placeholder.jpg'
---
`, file.Title, firstLine, formattedDate, fileID)

		path := "./app/src/content/blog"
		filePath := filepath.Join(path, fileID+".mdx")
		err = os.WriteFile(filePath, []byte(toAppend+md), 0644)
		if err != nil {
			return err
		}

		path = "./app/public/images"
		filePath = filepath.Join(path, fileID+"-placeholder.jpg")

		// get first image in google docs, it is the first <img tag in the html
		start := strings.Index(html, "<img")
		if start != -1 {
			start = strings.Index(html[start:], "src=\"") + start + 5
			end := strings.Index(html[start:], "\"") + start
			imageURL := html[start:end]
			fmt.Println(imageURL)
			err = downloadFile(filePath, imageURL)
			if err != nil {
				return err
			}
		} else {
			// blog-placeholder-about.jpg duplicate this file to the public folder
			file, err := os.Open("./app/public/blog-placeholder-about.jpg")
			if err != nil {
				return err
			}
			defer file.Close()
			err = createBlogFile(filePath, file)
			if err != nil {
				return err
			}

		}

	}
	posts = newPosts
	err = Build()
	if err != nil {
		return err
	}
	// copy the ./app/dist to ./web
	newDir := "./app/dist"
	err = CopyDir(newDir, "./web2")
	if err != nil {
		return err
	}

	err = os.Rename("./web", "./web3")
	if err != nil {
		return err
	}

	err = os.Rename("./web2", "./web")
	if err != nil {
		return err
	}

	err = os.RemoveAll("./web3")
	if err != nil {
		return err
	}

	lastChanged = newestModified

	// write it to files
	err = os.WriteFile("./files/lastmodified.txt", []byte(lastChanged), 0644)
	if err != nil {
		return err
	}
	return nil
}

func CopyFile(src, dst string) error {
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

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
func CopyDir(src, dst string) error {
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
			err := CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Copy file
			err := CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dest string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Skip copying if the source is not a file
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	// Copy file content
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// updateFiles copies new or modified files from newDir to oldDir.
func updateFiles(newDir, oldDir string) error {
	return filepath.Walk(newDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(newDir, path)
		if err != nil {
			return err
		}

		oldPath := filepath.Join(oldDir, relPath)

		// If the file doesn't exist in oldDir or is modified, copy it over
		oldFileInfo, err := os.Stat(oldPath)
		if os.IsNotExist(err) || oldFileInfo.ModTime().Before(info.ModTime()) {
			if info.Mode().IsRegular() {
				fmt.Printf("Copying %s to %s\n", path, oldPath)
				if err := copyFile(path, oldPath); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// deleteOldFiles removes files from oldDir if they don't exist in newDir.
func deleteOldFiles(newDir, oldDir string) error {
	return filepath.Walk(oldDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(oldDir, path)
		if err != nil {
			return err
		}

		newPath := filepath.Join(newDir, relPath)

		// If the file doesn't exist in newDir, delete it
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			if info.Mode().IsRegular() {
				fmt.Printf("Deleting %s\n", path)
				return os.Remove(path)
			}
		}

		return nil
	})
}

func downloadFile(filePath string, url string) error {
	// download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = createBlogFile(filePath, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func createBlogFile(filePath string, src io.Reader) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, src)
	if err != nil {
		return err
	}
	return nil
}

func removeFirstImage(md string) string {
	// remove the ![][image1]
	start := strings.Index(md, "![][image1]")
	if start == -1 {
		return md
	}
	end := start + len("![][image1]")
	md = md[:start] + md[end:]

	// remove [image1]: url
	start = strings.Index(md, "[image1]:")
	if start == -1 {
		return md
	}
	end = strings.Index(md[start:], "\n") + start
	md = md[:start] + md[end+1:]

	return md
}
