package cmd

import (
	"cmd/internal/command"
	"cmd/internal/config"
	"cmd/internal/database"
	"cmd/pkg/zap-logger"
	"os"
)

// Execute will setup database, configurations and cobra
func Execute() {
	configuration := config.LoadConfiguration()
	db := database.Database{}.Setup(configuration.Migration)
	cmd := command.InitCommands(configuration, db)

	if err := cmd.Execute(); err != nil {
		zap_logger.GetLogger().Error(err.Error())
		os.Exit(0)
	}
}
