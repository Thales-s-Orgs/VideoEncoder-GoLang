package service_test

import (
	"log"
	"os"
	"testing"

	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/application/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env")
	}
}

func TestUpload(t *testing.T) {

	video, repo := prepare()
	outputBucket := os.Getenv("outputBucketName")
	videoPath := os.Getenv("localStoragePath") + "/" + video.ID

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download(outputBucket)
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewUploadService()
	videoUpload.OutputBucket = outputBucket
	videoUpload.VideoPath = videoPath

	done := make(chan string)
	go videoUpload.ProccessUpload(50, done)
	result := <-done

	require.Equal(t, result, "upload completed")

}
