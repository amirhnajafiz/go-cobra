package handler

import (
	"cmd/internal/middleware"
	"cmd/internal/models"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewTaskHandler(db *gorm.DB) middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		_ = json.NewDecoder(r.Body).Decode(&task)

		task.Status = "Started"

		db.Create(&task)
		taskID := strconv.FormatUint(uint64(task.ID), 10)
		response := "{ \"task_id\" : " + taskID + "}"
		_, _ = fmt.Fprintf(w, response)

		// Starts running CaptainCore command
		// TODO
		// go runCommand("captaincore "+task.Command, task)
	}
}
