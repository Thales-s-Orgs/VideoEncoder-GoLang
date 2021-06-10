package services

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/application/repositories"
	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n", string(out))
	}
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {

	client, err := storage.NewClient(context.Background())
	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(v.Video.Filepath)
	r, err := obj.NewReader(context.Background())
	if err != nil {
		return err
	}

	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("Video %v has been stored", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {

	rootFolderPath := os.Getenv("localStoragePath") + "/" + v.Video.ID
	source, target := rootFolderPath+".mp4", rootFolderPath+".frag"

	err := os.Mkdir(rootFolderPath, os.ModePerm)
	if err != nil {
		return err
	}

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	rootFolderPath := os.Getenv("localStoragePath") + "/" + v.Video.ID

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, rootFolderPath+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, rootFolderPath)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	id := v.Video.ID
	rootFolderPath := os.Getenv("localStoragePath") + "/" + id

	err := os.Remove(rootFolderPath + ".mp4")
	if err != nil {
		log.Println("error removing mp4 ", id, ".mp4")
		return err
	}

	err = os.Remove(rootFolderPath + ".frag")
	if err != nil {
		log.Println("error removing frag ", id, ".frag")
		return err
	}

	err = os.RemoveAll(rootFolderPath)
	if err != nil {
		log.Println("error removing mp4 ", id, ".mp4")
		return err
	}

	log.Println("files have been removed: ", id)

	return nil

}
