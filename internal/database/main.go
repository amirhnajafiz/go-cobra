package database

import (
	"cmd/internal/models"
	logger "cmd/pkg/zap-logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initialMigration() {
	// Migrate the project schema
	err := db.AutoMigrate(&models.Task{})

	if err != nil {
		logger.GetLogger().Error("auto migration fail!")
	}
}

func Setup(migrate bool) *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open("sql.db"), &gorm.Config{})

	if err != nil {
		logger.GetLogger().Fatal("database connection fail!")
	}

	if migrate {
		initialMigration()
		logger.GetLogger().Info("migration done successfully.")
	}

	return db
}
