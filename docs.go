package docs_blog_engine

import (
	"fmt"
	"net/http"
)

func getDoc(docID string, format string) (string, error) {
	var mimeType string
	switch format {
	case "html":
		mimeType = "text/html"
	case "md":
		mimeType = "text/markdown"
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	file, err := srv.Files.Export(docID, mimeType).Download()
	if err != nil {
		return "", err
	}
	defer file.Body.Close()

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

func getHTMLandMD(docID string) (string, string, error) {
	html, err := getDoc(docID, "html")
	if err != nil {
		return "", "", err
	}
	md, err := getDoc(docID, "md")
	if err != nil {
		return "", "", err
	}
	return html, md, nil
}

func downloadFileLocally(filePath string, url string) error {
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
