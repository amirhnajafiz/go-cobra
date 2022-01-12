package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Command  string `json:"command"`
	Status   string `json:"status"`
	Response string `json:"response"`
}

func (t Task) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Command, validation.Required, validation.Length(0, 255)),
		validation.Field(&t.Status, validation.Required, validation.Length(0, 255)),
		validation.Field(&t.Response, validation.Required, validation.Length(0, 255)),
	)
}
