package handler

import (
	"cmd/internal/middleware"
	"cmd/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func ViewTaskHandler(db *gorm.DB) middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var tasks models.Task
		db.Where("id = ?", id).Find(&tasks)
		fmt.Println("{}", tasks)

		_ = json.NewEncoder(w).Encode(tasks)
	}
}
