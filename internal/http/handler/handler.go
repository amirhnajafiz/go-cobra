package handler

import (
	"github.com/amirhnajafiz/go-cobra/pkg/runner"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Runner runner.Runner
}
