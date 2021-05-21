package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	ID               string    `json:"encoded_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status" valid:"notnull"`
	Error            string    `valid:"-"`
	Video            *Video    `valid:"-"`
	VideoID          string    `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdatedAt        time.Time `json:"updated_at" valid:"-"`
}

func NewJob(path, status string, video *Video) (*Job, error) {

	job := Job{
		OutputBucketPath: path,
		Video:            video,
		Status:           status,
	}

	job.prepare()

	err := job.Validate()

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (j *Job) prepare() {

	j.ID = uuid.NewV4().String()
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()

}

func (j *Job) Validate() error {
	_, err := govalidator.ValidateStruct(j)

	if err != nil {
		return err
	}

	return nil
}
