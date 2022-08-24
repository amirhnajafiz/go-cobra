package http

import "github.com/amirhnajafiz/go-cobra/internal/models"

func Validate(task models.Task) error {
	return task.Validate()
}
