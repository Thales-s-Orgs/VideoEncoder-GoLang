package repositories_tests

import (
	"testing"
	"time"

	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/domain"
	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/framework/database"
	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/repositories/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoRepositoryDbInsert(t *testing.T) {
	db := database.NewDBTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.Filepath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{DB: db}
	repo.Insert(video)

	v, err := repo.Find(video.ID)

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)
}
