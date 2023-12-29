package drive

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type GoogleDrive struct {
	ClientEmail string         `json:"user_email"`
	service     *drive.Service `json:""`
}

func NewDriveService(ctx context.Context, credentialFilePath string) (*GoogleDrive, error) {
	srv, err := drive.NewService(ctx, option.WithCredentialsFile("private.json"),
		option.WithScopes(drive.DriveFileScope, drive.DriveScope))
	if err != nil {
		return nil, fmt.Errorf("failed to start drive service, error: %v", err)
	}
	return &GoogleDrive{
		ClientEmail: "",
		service:     srv,
	}, nil
}

type File struct {
	// contains information and you can config file permission and more
	driveFile *drive.File
	// contains file content that fit os specific
	File *os.File
	// folder is drive folder
	Folder []string
}

type UploadFileRequest struct {
	Files []File
}

func (gd *GoogleDrive) UploadFile(ctx context.Context, req *UploadFileRequest) ([]string, error) {
	respIDs := make([]string, len(req.Files))
	for i, f := range req.Files {
		// list file info
		fileInf, err := f.File.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to stat file [%d], error: %v", i, err)
		}

		// set name file for drive
		f.driveFile.Name = fileInf.Name()

		// do upload file to google drive
		resp, err := gd.service.Files.
			Create(f.driveFile).
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

// TODO: should check again, haven't list able yet!
func (gd *GoogleDrive) GetListFolder(ctx context.Context) (*drive.DriveList, error) {
	driveList, err := gd.service.Drives.List().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to to get list drives, error: %v", err)
	}
	return driveList, nil
}

func (gd *GoogleDrive) Delete(ctx context.Context, id string) error {
	return gd.service.Files.Delete(id).Do()
}
