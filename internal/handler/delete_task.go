package handler

import (
	"cmd/internal/middleware"
	"cmd/internal/models"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func DeleteTaskHandler(db *gorm.DB) middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		command := vars["command"]

		var tasks models.Task
		db.Where("command = ?", command).Find(&tasks)
		db.Delete(&tasks)

		_, _ = fmt.Fprintf(w, "Successfully Deleted Task")
	}
}
