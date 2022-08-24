package database

import (
	"github.com/amirhnajafiz/go-cobra/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbFileName = "sql.db"
)

// Migrate
// database schema.
func Migrate() error {
	conn, err := Connect()
	if err != nil {
		return err
	}

	// Migrate the project schema
	if er := conn.AutoMigrate(&models.Task{}); er != nil {
		return er
	}

	return nil
}

// Connect
// open connection to database.
func Connect() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
}
