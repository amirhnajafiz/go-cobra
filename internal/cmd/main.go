package cmd

import (
	"cmd/internal/command"
	"cmd/internal/config"
	"cmd/internal/database"
	"cmd/pkg/logger"
	"os"
)

// Execute will setup database, configurations and cobra
func Execute() {
	configuration := config.LoadConfiguration()
	log := logger.GetLogger()
	db := database.Database{
		Logger: log.Named("database"),
	}.Setup(configuration.Migration)
	cmd := command.Commander{
		DB:            db.DB,
		Configuration: configuration,
		Logger:        log.Named("commander"),
	}.InitCommands()

	if err := cmd.Execute(); err != nil {
		logger.GetLogger().Error(err.Error())
		os.Exit(0)
	}
}
