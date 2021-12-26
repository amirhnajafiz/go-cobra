package handler

import (
	"cmd/internal/middleware"
	"cmd/internal/models"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func UpdateTaskHandler(db *gorm.DB) middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		command := vars["command"]

		var tasks models.Task
		db.Where("command = ?", command).Find(&tasks)

		tasks.Command = command

		db.Save(&tasks)
		_, _ = fmt.Fprintf(w, "Successfully Updated Task")
	}
}
