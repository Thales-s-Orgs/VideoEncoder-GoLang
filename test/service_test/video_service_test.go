package service_test

import (
	"log"
	"testing"
	"time"

	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/application/repositories"
	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/application/services"
	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/domain"
	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/framework/database"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDBTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.Filepath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{DB: db}
	repo.Insert(video)

	return video, repo
}

func TestDownload(t *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("encodetest")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
	require.Nil(t, err)

}
