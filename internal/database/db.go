package database

import (
	"cmd/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (d Database) initialMigration() {
	// Migrate the project schema
	err := d.DB.AutoMigrate(&models.Task{})

	if err != nil {
		d.Logger.Error("auto migration fail!")
	}
}

func (d Database) Setup(migrate bool) Database {
	var err error
	d.DB, err = gorm.Open(sqlite.Open("sql.db"), &gorm.Config{})

	if err != nil {
		d.Logger.Fatal("database connection fail!")
	}

	if migrate {
		d.initialMigration()
		d.Logger.Info("migration done successfully.")
	}

	return d
}
