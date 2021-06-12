package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

	if _, err = io.Copy(wc, f); err != nil {
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

func (us *UploadService) ProccessUpload(concurrency int, done chan string) error {

	in := make(chan int, runtime.NumCPU())
	returnChan := make(chan string)

	err := us.LoadPaths()
	if err != nil {
		return err
	}

	client, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for proccess := 0; proccess < concurrency; proccess++ {
		go us.UploadWorker(in, returnChan, client, ctx)
	}

	go func() {
		for i := 0; i < len(us.Paths); i++ {
			in <- i
		}
		close(in)
	}()

	for v := range returnChan {
		if v != "" {
			done <- v
			break
		}
	}

	return nil
}

func (us *UploadService) UploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {
	for v := range in {

		err := us.UploadObject(us.Paths[v], uploadClient, ctx)
		if err != nil {
			us.Errors = append(us.Errors, us.Paths[v])
			log.Printf("Error during upload: %v. Error %v", us.Paths[v], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- "upload completed"
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	return client, ctx, nil
}
