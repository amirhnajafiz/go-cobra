package cmd

import (
	"cmd/config"
	"cmd/internal/command"
	"cmd/internal/database"
	"cmd/pkg/zap-logger"
	"os"
)

func Exec() {
	db := database.Setup(true)
	configuration := config.GetConfig()
	cmd := command.InitCommands(configuration, db)

	if err := cmd.Execute(); err != nil {
		zap_logger.GetLogger().Error(err.Error())
		os.Exit(0)
	}
}
