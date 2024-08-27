package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/quocbang/multi-google-drive-upload/account"
	localDrive "github.com/quocbang/multi-google-drive-upload/drive"
	"google.golang.org/api/drive/v3"
)

type Env struct {
	APIKey string `envconfig:"API_KEY"`
}

func loadEnv() (*Env, error) {
	// load default .env file, ignore the error
	_ = godotenv.Load()

	env := new(Env)
	err := envconfig.Process("", env)
	if err != nil {
		return nil, fmt.Errorf("load config error: %v", err)
	}

	return env, nil
}

func main() {
	_, err := loadEnv()
	if err != nil {
		log.Fatalf("failed to load env, error: %v", err)
	}

	account, err := account.NewAccount()
	if err != nil {
		log.Fatal(err)
	}
	account.GetUserInfo("your_access_token")

	ctx := context.Background()
	service, err := localDrive.NewDriveService(ctx, "private.json")
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Open("./assets/video/mo_denvau.mp4")
	if err != nil {
		log.Fatal(err)
	}

	ids, err := service.UploadFile(ctx, &localDrive.UploadFileRequest{
		Files: []localDrive.File{
			{
				DriveFile: &drive.File{},
				Folder:    []string{"10REQ7-cIKN46ymMsEOZHt6HgQXC4K0Nu"},
				File:      f,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ids)
}
