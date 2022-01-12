package handler

import (
	"cmd/internal/middleware"
	"cmd/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func AllTasksHandler(db *gorm.DB) middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tasks []models.Task

		vars := mux.Vars(r)
		page, _ := strconv.Atoi(vars["page"])

		if page > 0 {
			offset := page * 10
			db.Offset(offset).Limit(10).Order("created_at desc").Find(&tasks)
		} else {
			db.Limit(10).Order("created_at desc").Find(&tasks)
		}

		_ = json.NewEncoder(w).Encode(tasks)
	}
}
