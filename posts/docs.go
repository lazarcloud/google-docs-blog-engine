package posts

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lazarcloud/google-docs-blog-engine/globals"
	"google.golang.org/api/drive/v2"
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

func getHTMLandMD(docID string) (html string, md string, err error) {
	html, err = getDoc(docID, "html")
	if err != nil {
		return "", "", err
	}
	md, err = getDoc(docID, "md")
	if err != nil {
		return "", "", err
	}
	return html, md, nil
}

func saveFile(filePath string, src io.Reader) error {
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

func downloadFileLocally(filePath string, url string) error {
	// download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = saveFile(filePath, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func removeFirstImage(md string) string {
	start := strings.Index(md, "![][image1]")
	if start == -1 {
		return md
	}
	end := start + len("![][image1]")
	md = md[:start] + md[end:]

	start = strings.Index(md, "[image1]:")
	if start == -1 {
		return md
	}
	end = strings.Index(md[start:], "\n") + start
	md = md[:start] + md[end+1:]

	return md
}

func getDescription(input string) (md string, description string) {
	firstLine := strings.Split(input, "\n")[0]
	if !strings.HasPrefix(firstLine, globals.DescriptionKeyword) {
		return input, globals.DefaultDescription
	}
	firstLine = strings.TrimPrefix(firstLine, globals.DescriptionKeyword)
	md = strings.Join(strings.Split(input, "\n")[1:], "\n")
	return md, firstLine
}

func savePicture(postImage string, html string) error {
	filePath := filepath.Join(globals.ImagesRoot, postImage)

	start := strings.Index(html, "<img")
	if start != -1 {
		start = strings.Index(html[start:], "src=\"") + start + 5
		end := strings.Index(html[start:], "\"") + start
		imageURL := html[start:end]
		fmt.Println(imageURL)
		err := downloadFileLocally(filePath, imageURL)
		return err
	}

	file, err := os.Open(globals.DefaultImagePath)
	if err != nil {
		return err
	}
	defer file.Close()
	err = saveFile(filePath, file)
	return err
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

func formatDate(input string) (string, error) {
	parsedDate, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return "", err
	}

	formattedDate := parsedDate.Format("02 Jan 2006")
	return formattedDate, nil
}
