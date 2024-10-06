package docs_blog_engine

import (
	"fmt"
	"os"
	"time"

	notrhttp "github.com/Notr-Dev/notr-http"
	notrhttp_middlewares "github.com/Notr-Dev/notr-http/middlewares"
	files "github.com/lazarcloud/google-docs-blog-engine/fs"
	"github.com/lazarcloud/google-docs-blog-engine/posts"
	"github.com/lazarcloud/google-docs-blog-engine/run"
)

func RunServer() error {
	err := files.EnsureFoldersExist([]string{"./files", "./web"})
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

	err = run.Install()
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
	err = posts.CheckIfConfigured()
	if err != nil {
		return err
	}

	_, err = os.Stat("./web/index.html")
	if err != nil {

		err = posts.GetPosts()
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
			posts.GetPosts()
			return nil
		},
	})

	err = server.Run()
	return err
}
