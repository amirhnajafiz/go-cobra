package migration

import (
	"github.com/amirhnajafiz/go-cobra/internal/database"
	"github.com/amirhnajafiz/go-cobra/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "migrate",
		Long:  "migrate database (sqlite3)",
		Run: func(_ *cobra.Command, _ []string) {
			main()
		},
	}
}

func main() {
	l := logger.New()

	if err := database.Migrate(); err != nil {
		l.Error("migration failed", zap.Error(err))

		return
	}

	l.Info("migrate successfully")
}
