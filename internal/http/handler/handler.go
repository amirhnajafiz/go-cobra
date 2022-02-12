package handler

import (
	"cmd/pkg/runner"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Runner runner.Runner
}
