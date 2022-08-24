package cmd

import (
	"os"

	"github.com/amirhnajafiz/go-cobra/internal/command"
	"github.com/amirhnajafiz/go-cobra/internal/config"
	"github.com/amirhnajafiz/go-cobra/internal/database"
	"github.com/amirhnajafiz/go-cobra/pkg/logger"
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
