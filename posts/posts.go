package posts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lazarcloud/google-docs-blog-engine/backup"
	files "github.com/lazarcloud/google-docs-blog-engine/fs"
	"github.com/lazarcloud/google-docs-blog-engine/globals"
	docs_blog_engine_run "github.com/lazarcloud/google-docs-blog-engine/run"
)

type Post struct {
	Title     string `json:"title"`
	ID        string `json:"id"`
	Date      string `json:"date"`
	ContentMD string `json:"content_md"`
}

func GetPosts(toWait time.Duration) error {
	newestModified, fileList, err := getLastModified()
	if err != nil {
		return err
	}

	if newestModified == lastChanged {
		fmt.Println("No changes in the folder")
		return nil
	}

	fmt.Println("Changes detected in the folder")

	time.Sleep(toWait)

	secondModified, fileList, err := getLastModified()
	if err != nil {
		return err
	}

	if secondModified != newestModified {
		return nil
	}

	err = files.ClearDirectories([]string{"./app/src/content/blog", "./app/public/images"})
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

		md, description := getDescription(md)

		formattedDate, err := formatDate(file.CreatedDate)
		if err != nil {
			return err
		}

		postImage := postID + "-placeholder.jpg"

		md, err = fixImages(md, html)
		if err != nil {
			return err
		}

		err = savePicture(postImage, html)
		if err != nil {
			return err
		}
		md = removeFirstImage(md)

		toAppend := fmt.Sprintf(`---
title: '%s'
description: '%s'
pubDate: '%s'
heroImage: '/images/%s-placeholder.jpg'
---
`, file.Title, description, formattedDate, postID)

		md = strings.ReplaceAll(md, "\\#", "#")
		md = strings.ReplaceAll(md, "\\.", ".")
		md = strings.ReplaceAll(md, "\\`", "`")
		md = strings.ReplaceAll(md, "\\-", "-")

		path := "./app/src/content/blog"
		filePath := filepath.Join(path, postID+".mdx")
		err = os.WriteFile(filePath, []byte(toAppend+md), 0644)
		if err != nil {
			return err
		}

	}
	err = docs_blog_engine_run.Build()
	if err != nil {
		return err
	}

	err = files.CopyDir(globals.BuildDir, globals.StaticDir+"_new")
	if err != nil {
		return err
	}

	err = os.Rename(globals.StaticDir, globals.StaticDir+"_redacted")
	if err != nil {
		return err
	}

	err = os.Rename(globals.StaticDir+"_new", globals.StaticDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(globals.StaticDir + "_redacted")
	if err != nil {
		return err
	}

	err = backup.CreateBackup()
	if err != nil {
		return err
	}

	lastChanged = newestModified

	return nil
}
