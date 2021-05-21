package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Video struct {
	ID         string    `json:"encoded_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	ResourceID string    `json:"resource_id" valid:"notnull" gorm:"type:varchar(255)"`
	Filepath   string    `json:"filepath" valid:"notnull" gorm:"type:varchar(255)"`
	CreatedAt  time.Time `valid:"-"`
	Jobs       []*Job    `json:"jobs" valid:"notnull" gorm:"ForeignKey:VideoID"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewVideo() *Video {
	return &Video{}
}

func (v *Video) Validate() error {
	_, err := govalidator.ValidateStruct(v)

	if err != nil {
		return err
	}

	return nil
}
