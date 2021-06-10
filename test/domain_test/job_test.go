package domain_test

import (
	"testing"
	"time"

	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.Filepath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("path", "Started", video)
	require.NotNil(t, job)
	require.Nil(t, err)
}
