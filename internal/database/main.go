package database

import (
	"cmd/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initialMigration() {
	// Migrate the project schema
	_ = db.AutoMigrate(&models.Task{})
}

func Setup(migrate bool) *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open("sql.db"), &gorm.Config{})

	if err != nil {
		panic("database connection: faild")
	}

	if migrate {
		initialMigration()
	}

	return db
}
