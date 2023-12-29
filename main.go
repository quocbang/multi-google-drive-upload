package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
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

type GoogleDrive struct {
	ClientEmail string         `json:"user_email"`
	Service     *drive.Service `json:""`
}

func InitService(ctx context.Context) (*drive.Service, error) {
	srv, err := drive.NewService(ctx, option.WithCredentialsFile("private.json"),
		option.WithScopes(drive.DriveFileScope, drive.DriveScope))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func main() {
	_, err := loadEnv()
	if err != nil {
		log.Fatalf("failed to load env, error: %v", err)
	}

	ctx := context.Background()
	srv, err := InitService(ctx)
	if err != nil {
		log.Fatalf("failed to start drive service, error: %v", err)
	}

	gd := &GoogleDrive{
		ClientEmail: "", // TODO: should get in private file
		Service:     srv,
	}

	driveList, err := gd.getListFolder(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(driveList)

	os.Exit(0)
	// file path
	filePath := "./assets/video/mo_denvau.mp4"
	// baseMineType := "multipart/file"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file, error: %v", err)
	}
	defer file.Close()

	f := &drive.File{
		Parents: []string{"10REQ7-cIKN46ymMsEOZHt6HgQXC4K0Nu"},
		// WritersCanShare: true,
		// Shared:          true,
		// Permissions: []*drive.Permission{
		// 	{
		// 		Role:   "writer",
		// 		Type:   "domain",
		// 		Domain: "teqnological.asia",
		// 	},
		// },
	}

	uploadFileRequest := &UploadFileRequest{
		Files: []File{
			{
				DriveFile: f,
				File:      file,
			},
		},
	}
	respIDs, err := gd.uploadFile(ctx, uploadFileRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(respIDs)
}

type File struct {
	// contains information and you can config file permission and more
	DriveFile *drive.File
	// contains file content that fit os specific
	File *os.File
}

type UploadFileRequest struct {
	Files []File
}

func (gd *GoogleDrive) uploadFile(ctx context.Context, req *UploadFileRequest) ([]string, error) {
	respIDs := make([]string, len(req.Files))
	for i, f := range req.Files {
		// list file info
		fileInf, err := f.File.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to stat file [%d], error: %v", i, err)
		}

		// set name file for drive
		f.DriveFile.Name = fileInf.Name()

		// do upload file to google drive
		resp, err := gd.Service.Files.
			Create(f.DriveFile).
			Media(f.File, googleapi.ChunkSize(int(fileInf.Size()))).
			ProgressUpdater(func(current, total int64) {
				fmt.Println(current, total)
			}).
			Do()
		if err != nil {
			return nil, fmt.Errorf("failed to upload file [%d], error: %v", i, err)
		}
		respIDs[i] = resp.Id
	}
	return respIDs, nil
}

func (gd *GoogleDrive) getListFolder(ctx context.Context) (*drive.DriveList, error) {
	driveList, err := gd.Service.Drives.List().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to to get list drives, error: %v", err)
	}
	return driveList, nil
}

func (gd *GoogleDrive) delete(ctx context.Context, id string) error {
	return gd.Service.Files.Delete(id).Do()
}
