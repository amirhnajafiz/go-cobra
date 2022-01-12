package request

import "cmd/internal/models"

func Validate(task models.Task) error {
	return task.Validate()
}
