package command

import (
	"github.com/amirhnajafiz/go-cobra/internal/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Commander struct {
	Configuration config.Config
	DB            *gorm.DB
	Logger        *zap.Logger
}
