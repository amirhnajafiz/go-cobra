package migration

import (
	"github.com/amirhnajafiz/go-cobra/internal/database"
	"github.com/amirhnajafiz/go-cobra/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	l := logger.New()

	if err := database.Migrate(); err != nil {
		l.Error("migration failed", zap.Error(err))
	}
}