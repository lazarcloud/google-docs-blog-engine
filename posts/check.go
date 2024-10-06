package posts

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

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
