package docs_blog_engine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/api/drive/v2"
)

type Post struct {
	Title     string `json:"title"`
	ID        string `json:"id"`
	Date      string `json:"date"`
	ContentMD string `json:"content_md"`
}

func getLastModified() (string, *drive.FileList, error) {
	folder, err := srv.Files.Get(string(folderID)).Do()
	if err != nil {
		return "", &drive.FileList{}, err
	}

	dates := []string{folder.ModifiedDate}

	fileList, err := srv.Files.List().Q(fmt.Sprintf("'%s' in parents", folderID)).Do()
	if err != nil {
		return "", &drive.FileList{}, err
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
	return newestModified, fileList, nil
}

func getPosts() error {
	newestModified, fileList, err := getLastModified()
	if err != nil {
		return err
	}

	if newestModified == lastChanged {
		fmt.Println("No changes in the folder")
		return nil
	}

	fmt.Println("Changes detected in the folder")

	err = clearDirectories([]string{"./app/src/content/blog", "./app/public/images"})
	if err != nil {
		return err
	}

	for _, file := range fileList.Items {
		fmt.Printf("Name: %s, ID: %s\n", file.Title, file.Id)
		html, md, err := getHTMLandMD(file.Id)
		if err != nil {
			return err
		}

		postID := strings.ToLower(strings.ReplaceAll(file.Title, " ", "-"))

		firstLine := strings.Split(md, "\n")[0]
		if !strings.HasPrefix(firstLine, "DESCRIPTION ") {
			firstLine = "No description"
		} else {
			firstLine = strings.TrimPrefix(firstLine, "DESCRIPTION ")
			md = strings.Join(strings.Split(md, "\n")[1:], "\n")
		}

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
`, file.Title, firstLine, formattedDate, postID)

		path := "./app/src/content/blog"
		filePath := filepath.Join(path, postID+".mdx")
		err = os.WriteFile(filePath, []byte(toAppend+md), 0644)
		if err != nil {
			return err
		}

		path = "./app/public/images"
		filePath = filepath.Join(path, postID+"-placeholder.jpg")

		// get first image in google docs, it is the first <img tag in the html
		start := strings.Index(html, "<img")
		if start != -1 {
			start = strings.Index(html[start:], "src=\"") + start + 5
			end := strings.Index(html[start:], "\"") + start
			imageURL := html[start:end]
			fmt.Println(imageURL)
			err = downloadFileLocally(filePath, imageURL)
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
	err = Build()
	if err != nil {
		return err
	}
	// copy the ./app/dist to ./web
	newDir := "./app/dist"
	err = copyDir(newDir, "./web2")
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

	err = os.WriteFile("./files/lastmodified.txt", []byte(lastChanged), 0644)
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
