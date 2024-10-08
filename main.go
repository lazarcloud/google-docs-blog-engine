package docs_blog_engine

import (
	"time"

	notrhttp "github.com/Notr-Dev/notr-http"
	notrhttp_middlewares "github.com/Notr-Dev/notr-http/middlewares"
	files "github.com/lazarcloud/google-docs-blog-engine/fs"
	"github.com/lazarcloud/google-docs-blog-engine/globals"
	"github.com/lazarcloud/google-docs-blog-engine/posts"
	"github.com/lazarcloud/google-docs-blog-engine/run"
)

func RunServer() error {
	err := files.EnsureFoldersExist([]string{"./web"})
	if err != nil {
		return err
	}

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

	err = posts.GetPosts()
	if err != nil {
		return err
	}

	server.ServeStaticWebsite("/", globals.StaticDir)

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
