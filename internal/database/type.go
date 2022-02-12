package database

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Database struct {
	DB     *gorm.DB
	Logger *zap.Logger
}
