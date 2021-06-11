package services

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

type UploadService struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (us *UploadService) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {

	path := strings.Split(objectPath, os.Getenv("localStoragePath")+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(us.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err := io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (us *UploadService) LoadPaths() error {

	err := filepath.Walk(us.VideoPath, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			us.Paths = append(us.Paths, path)
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil
}

func (us *UploadService) getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
