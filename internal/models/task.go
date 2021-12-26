package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Command  string `json:"command"`
	Status   string `json:"status"`
	Response string `json:"response"`
}
