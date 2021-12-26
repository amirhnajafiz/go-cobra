package handler

import (
	"cmd/internal/models"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

func NewRunHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			_, _ = fmt.Fprintf(w, "404 wrong input")
			return
		}

		task.Status = "Started"

		db.Create(&task)

		// Starts running CaptainCore command
		response := "command" // runCommand("captaincore "+task.Command, task)
		_, _ = fmt.Fprintf(w, response)

	}
}
