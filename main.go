package docs_blog_engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	notrhttp "github.com/Notr-Dev/notr-http"
	notrhttp_middlewares "github.com/Notr-Dev/notr-http/middlewares"
	docs_blog_engine_run "github.com/lazarcloud/google-docs-blog-engine/run"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

var srv *drive.Service
var folderID []byte
var lastChanged = ""

func CheckIfConfigured() error {
	_, err := os.Stat("./files/lastmodified.txt")
	if err != nil {
		_, err = os.Create("./files/lastmodified.txt")
		if err != nil {
			return err
		}
		lastChanged = ""
	} else {
		file, err := os.ReadFile("./files/lastmodified.txt")
		if err != nil {
			return err
		}
		lastChanged = string(file)
	}
	ctx := context.Background()
	if os.Getenv("GOOGLE_CREDENTIALS") == "" {
		return fmt.Errorf("GOOGLE_CREDENTIALS env variable is not set")
	}
	config, err := google.JWTConfigFromJSON([]byte(os.Getenv("GOOGLE_CREDENTIALS")), drive.DriveScope)
	if err != nil {
		return err
	}
	client := config.Client(ctx)
	srv, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	if os.Getenv("GOOGLE_FOLDER_ID") == "" {
		return fmt.Errorf("GOOGLE_FOLDER_ID env variable is not set")
	}
	folderID = []byte(os.Getenv("GOOGLE_FOLDER_ID"))
	_, err = srv.Files.List().Q(fmt.Sprintf("'%s' in parents", folderID)).Do()
	if err != nil {
		if strings.Contains(err.Error(), "googleapi: Error 403: Google Drive API has not been used in project") {
			return errors.New("Drive API has not been enabled")
		} else {
			return err
		}
	}
	return nil
}

func RunServer() error {
	err := ensureFoldersExist([]string{"./files", "./web"})
	if err != nil {
		return err
	}

	err = os.WriteFile("./files/credentials.json", []byte(os.Getenv("GOOGLE_CREDENTIALS")), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile("./files/folder.txt", []byte(os.Getenv("GOOGLE_FOLDER_ID")), os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("Managed to write credentials")

	err = docs_blog_engine_run.Install()
	if err != nil {
		return err
	}

	server := notrhttp.NewServer(
		notrhttp.Server{
			Name:    "Google Docs Blog",
			Port:    ":8080",
			Version: "0.0.1",
		},
	)
	err = CheckIfConfigured()
	if err != nil {
		return err
	}

	_, err = os.Stat("./web/index.html")
	if err != nil {

		err = getPosts()
		if err != nil {
			return err
		}
	}

	server.ServeStaticWebsite("/", "./web")

	server.RegisterMiddleware(notrhttp_middlewares.AllowALlOrigins)

	server.RegisterJob(notrhttp.Job{
		Name:     "Get Posts",
		Interval: time.Second * 30,
		Job: func() error {
			getPosts()
			return nil
		},
	})

	err = server.Run()
	return err
}
