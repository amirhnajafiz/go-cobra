package handler

import (
	http2 "cmd/internal/http"
	"cmd/internal/middleware"
	"cmd/internal/models"
	commander "cmd/pkg/command-runner"
	"encoding/json"
	"fmt"
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

func NewRunHandler(db *gorm.DB) middleware.HttpHandlerFunc {
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
		response := commander.RunCommand("captain-core "+task.Command, task, db)
		_, _ = fmt.Fprintf(w, response)

	}
}

func NewTaskHandler(db *gorm.DB) middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		_ = json.NewDecoder(r.Body).Decode(&task)

		err := http2.Validate(task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, err.Error())
		}

		task.Status = "Started"

		db.Create(&task)
		taskID := strconv.FormatUint(uint64(task.ID), 10)
		response := "{ \"task_id\" : " + taskID + "}"
		_, _ = fmt.Fprintf(w, response)

		// Starts running Captain-Core command
		go commander.RunCommand("captain-core "+task.Command, task, db)
	}
}

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
