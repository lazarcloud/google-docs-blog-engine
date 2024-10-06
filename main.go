package docs_blog_engine

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	notrhttp "github.com/Notr-Dev/notr-http"
	notrhttp_middlewares "github.com/Notr-Dev/notr-http/middlewares"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

var srv *drive.Service
var folderID []byte
var lastChanged = ""

func CheckIfConfigured() (string, error, string) {
	// create the files directory
	// check if file ./files/lastmodified.txt exists
	_, err := os.Stat("./files/lastmodified.txt")
	if err != nil {
		// create it and dont write anything in it
		_, err = os.Create("./files/lastmodified.txt")
		if err != nil {
			return configurationsStages[0], err, ""
		}
		lastChanged = ""
	} else {
		file, err := os.ReadFile("./files/lastmodified.txt")
		if err != nil {
			return configurationsStages[0], err, ""
		}
		lastChanged = string(file)
	}
	ctx := context.Background()
	creds, err := os.ReadFile("./files/credentials.json")
	if err != nil {
		return configurationsStages[0], err, ""
	}
	config, err := google.JWTConfigFromJSON(creds, drive.DriveScope)
	if err != nil {
		return configurationsStages[0], err, ""
	}
	client := config.Client(ctx)
	srv, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return configurationsStages[1], err, ""
	}
	folderID, err = os.ReadFile("./files/folder.txt")
	if err != nil {
		return configurationsStages[1], err, ""
	}
	// List all files in the specified folder
	fileList, err := srv.Files.List().Q(fmt.Sprintf("'%s' in parents", folderID)).Do()
	if err != nil {
		// log.Fatalf("Unable to retrieve files in folder: %v", err.Error())
		if strings.Contains(err.Error(), "googleapi: Error 403: Google Drive API has not been used in project") {
			// found beetween Enable it by visiting and  then retry
			start := strings.Index(err.Error(), "https://")
			end := strings.Index(err.Error()[start:], " ") + start
			url := err.Error()[start:end]
			return configurationsStages[2], err, url
		} else {
			return configurationsStages[1], err, ""
		}
	}
	for _, file := range fileList.Items {
		fmt.Printf("Name: %s, ID: %s\n", file.Title, file.Id)
		fmt.Println(file.CreatedDate)
		fmt.Println(file.MimeType)

	}
	return configurationsStages[3], nil, ""
}

var configurationsStages = []string{
	"nothing",
	"credentials",
	"folder link",
	"drive api enabled",
}

var status string = "bad"

func firstCheck() error {
	progress, err, extra := CheckIfConfigured()
	fmt.Println(progress, err, extra)
	if progress != configurationsStages[len(configurationsStages)-1] {
		return nil
	}
	status = "ok"

	return nil
}

func RunServer() error {

	err := os.MkdirAll("./files", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll("./web", os.ModePerm)
	if err != nil {
		return err
	}

	// write to credentials.json the env variable GOOGLE_CREDENTIALS
	err = os.WriteFile("./files/credentials.json", []byte(os.Getenv("GOOGLE_CREDENTIALS")), os.ModePerm)
	if err != nil {
		return err
	}
	// write to folder.txt the env variable GOOGLE_FOLDER
	err = os.WriteFile("./files/folder.txt", []byte(os.Getenv("GOOGLE_FOLDER_ID")), os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("Written creds")

	err = Install()
	if err != nil {
		return err
	}

	// time.Sleep(time.Second * 100)

	server := notrhttp.NewServer(
		notrhttp.Server{
			Name:    "Google Docs Blog",
			Port:    ":8080",
			Version: "0.0.1",
		},
	)
	err = firstCheck()
	if err != nil {
		return err
	}
	// check if the web folder exists and has an index.html in it
	_, err = os.Stat("./web/index.html")
	if err != nil {

		err = getPosts()
		if err != nil {
			return err
		}
	}

	server.Get("/posts", func(rw notrhttp.Writer, r *notrhttp.Request) {
		rw.RespondWithSuccess(posts)
	})

	server.ServeStaticWebsite("/", "./web")

	server.RegisterMiddleware(notrhttp_middlewares.AllowALlOrigins)

	server.RegisterJob(notrhttp.Job{
		Name:     "Get Posts",
		Interval: time.Second * 10,
		Job: func() error {
			if status == "bad" {
				return nil
			}
			fmt.Println("Getting posts")
			getPosts()
			return nil
		},
	})

	err = server.Run()

	if err != nil {
		return err
	}
	return nil
}
