package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/amirhnajafiz/go-cobra/internal/http/middleware"
	"github.com/amirhnajafiz/go-cobra/internal/models"
	"github.com/amirhnajafiz/go-cobra/internal/runner"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Runner runner.Runner
}

func (h Handler) AllTasks() middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tasks []models.Task

		vars := mux.Vars(r)
		page, _ := strconv.Atoi(vars["page"])

		if page > 0 {
			offset := page * 10

			h.DB.Offset(offset).Limit(10).Order("created_at desc").Find(&tasks)
		} else {
			h.DB.Limit(10).Order("created_at desc").Find(&tasks)
		}

		_ = json.NewEncoder(w).Encode(tasks)
	}
}

func (h Handler) DeleteTask() middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		command := vars["command"]

		var tasks models.Task

		h.DB.Where("command = ?", command).Find(&tasks)
		h.DB.Delete(&tasks)

		_, _ = fmt.Fprintf(w, "Successfully Deleted Task")
	}
}

func (h Handler) NewRun() middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task

		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			_, _ = fmt.Fprintf(w, "404 wrong input")

			return
		}

		task.Status = "Started"

		h.DB.Create(&task)

		// Starts running CaptainCore command
		response := h.Runner.RunCommand("captain-core "+task.Command, task)

		_, _ = fmt.Fprintf(w, response)

	}
}

func (h Handler) NewTask() middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task

		_ = json.NewDecoder(r.Body).Decode(&task)

		err := task.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			_, _ = fmt.Fprint(w, err.Error())

			return
		}

		task.Status = "Started"

		h.DB.Create(&task)
		taskID := strconv.FormatUint(uint64(task.ID), 10)
		response := "{ \"task_id\" : " + taskID + "}"
		_, _ = fmt.Fprintf(w, response)

		// Starts running Captain-Core command
		go h.Runner.RunCommand("captain-core "+task.Command, task)
	}
}

func (h Handler) UpdateTask() middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		command := vars["command"]

		var tasks models.Task

		h.DB.Where("command = ?", command).Find(&tasks)

		tasks.Command = command

		h.DB.Save(&tasks)

		_, _ = fmt.Fprintf(w, "Successfully Updated Task")
	}
}

func (h Handler) ViewTask() middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var tasks models.Task

		h.DB.Where("id = ?", id).Find(&tasks)

		fmt.Println("{}", tasks)

		_ = json.NewEncoder(w).Encode(tasks)
	}
}
