package command

import (
	"cmd/internal/config"
	"gorm.io/gorm"
)

type Commander struct {
	Configuration config.Config
	DB            *gorm.DB
}
