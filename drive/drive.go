package drive

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/panjf2000/ants/v2"
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
	DriveFile *drive.File
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
	chIDs := make(chan string, 10)
	chErr := make(chan error)

	pool, err := ants.NewPoolWithFunc(len(req.Files), func(i interface{}) {
		f, ok := i.(File)
		if !ok {
			chErr <- fmt.Errorf("wrong type")
			return
		}
		// list file info
		fileInf, err := f.File.Stat()
		if err != nil {
			chErr <- fmt.Errorf("failed to stat file [%d], error: %v", i, err)
			return
		}

		// set name file for drive
		f.DriveFile.Name = fileInf.Name()

		// set folder where file is store
		f.DriveFile.Parents = f.Folder

		// do upload file to google drive
		resp, err := gd.service.Files.
			Create(f.DriveFile).
			Media(f.File, googleapi.ChunkSize(int(fileInf.Size()))).
			ProgressUpdater(func(current, total int64) {
				fmt.Println(current, total)
			}).
			Do()
		if err != nil {
			chErr <- fmt.Errorf("failed to upload file [%d], error: %v", i, err)
			return
		}
		chIDs <- resp.Id
	})
	defer pool.Release()
	if err != nil {
		return nil, fmt.Errorf("failed to create a new pool, error: %v", err)
	}

	for _, file := range req.Files {
		err := pool.Invoke(file)
		if err != nil {
			return nil, fmt.Errorf("failed to invoke, error: %v", err)
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for err := range chErr {
			if err != nil {
				close(chErr)
				close(chIDs)
				log.Fatal(err)
			}
		}
	}()
	go func(ids []string) {
		defer wg.Done()
		count := 0
		for id := range chIDs {
			count++
			respIDs = append(respIDs, id)
			if count == len(req.Files) {
				close(chErr)
				close(chIDs)
				return
			}
			log.Println(respIDs)
		}
	}(respIDs)
	wg.Wait()
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

// options with unexposed for control option that user input.
type options struct {
	folderID string
}

// Option is definition optional function.
type Option func(*options)

func ListFileWithFolder(folderID string) Option {
	return func(o *options) {
		o.folderID = folderID
	}
}

// parseOptions is definition optional function to parse value into
// option struct.
func parseOptions(opt ...Option) options {
	opts := &options{}

	for _, o := range opt {
		o(opts)
	}

	return *opts
}

// ListFiles is list file ids
func (gd *GoogleDrive) ListFiles(ctx context.Context, opt ...Option) ([]*drive.File, error) {
	// opts := parseOptions(opt...)
	result, err := gd.service.Files.List().Do()
	if err != nil {
		return nil, err
	}

	return result.Files, nil
}

func (gd *GoogleDrive) EmptyTrash(ctx context.Context) error {
	return gd.service.Files.EmptyTrash().Do()
}
