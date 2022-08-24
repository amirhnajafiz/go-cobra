package handler

import (
	"encoding/json"
	"fmt"
	http2 "github.com/amirhnajafiz/go-cobra/internal/http"
	"github.com/amirhnajafiz/go-cobra/internal/middleware"
	"github.com/amirhnajafiz/go-cobra/internal/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h Handler) AllTasksHandler() middleware.HttpHandlerFunc {
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

func (h Handler) DeleteTaskHandler() middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		command := vars["command"]

		var tasks models.Task
		h.DB.Where("command = ?", command).Find(&tasks)
		h.DB.Delete(&tasks)

		_, _ = fmt.Fprintf(w, "Successfully Deleted Task")
	}
}

func (h Handler) NewRunHandler() middleware.HttpHandlerFunc {
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

func (h Handler) NewTaskHandler() middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		_ = json.NewDecoder(r.Body).Decode(&task)

		err := http2.Validate(task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, err.Error())
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

func (h Handler) UpdateTaskHandler() middleware.HttpHandlerFunc {
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

func (h Handler) ViewTaskHandler() middleware.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var tasks models.Task
		h.DB.Where("id = ?", id).Find(&tasks)
		fmt.Println("{}", tasks)

		_ = json.NewEncoder(w).Encode(tasks)
	}
}
