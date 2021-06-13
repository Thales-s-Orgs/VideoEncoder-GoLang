package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/application/repositories"
	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

var (
	inputBucket  = os.Getenv("inputBucketName")
	outputBucket = os.Getenv("outputBucketName")
	videoPath    = os.Getenv("localStoragePath") + "/"
	concurrency  = os.Getenv("CONCURRENCY_UPLOAD")
	downloading  = "DOWNLOADING"
	fragmenting  = "FRAGMENTING"
	encoding     = "ENCODING"
	uploading    = "UPLOADING"
	finishing    = "FINISHING"
	completed    = "COMPLETED"
)

func (j *JobService) Start() error {
	err := j.ChangeJobStatus(downloading)
	if err != nil {
		return j.FailJob(err)
	}

	err = j.VideoService.Download(inputBucket)
	if err != nil {
		return j.FailJob(err)
	}

	err = j.ChangeJobStatus(fragmenting)
	if err != nil {
		return j.FailJob(err)
	}

	err = j.VideoService.Fragment()
	if err != nil {
		return j.FailJob(err)
	}

	err = j.ChangeJobStatus(encoding)
	if err != nil {
		return j.FailJob(err)
	}

	err = j.VideoService.Encode()
	if err != nil {
		return j.FailJob(err)
	}

	err = j.performUpload()
	if err != nil {
		return j.FailJob(err)
	}

	err = j.ChangeJobStatus(finishing)
	if err != nil {
		return j.FailJob(err)
	}

	err = j.VideoService.Finish()
	if err != nil {
		return j.FailJob(err)
	}

	err = j.ChangeJobStatus(completed)
	if err != nil {
		return j.FailJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {
	err := j.ChangeJobStatus(uploading)
	if err != nil {
		return j.FailJob(err)
	}

	uploadService := NewUploadService()
	uploadService.OutputBucket = outputBucket
	uploadService.VideoPath = videoPath + j.VideoService.Video.ID
	concurrencyInt, _ := strconv.Atoi(concurrency)
	done := make(chan string)

	go uploadService.ProccessUpload(concurrencyInt, done)
	uploadResult := <-done

	if uploadResult != "upload completed" {
		return j.FailJob(errors.New(uploadResult))
	}

	return nil
}

func (j *JobService) ChangeJobStatus(status string) error {
	j.Job.Status = status

	_, err := j.JobRepository.Update(j.Job)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobService) FailJob(erro error) error {
	j.Job.Status = "FAILED"
	j.Job.Error = erro.Error()

	_, err := j.JobRepository.Update(j.Job)
	if err != nil {
		return err
	}

	return nil
}
