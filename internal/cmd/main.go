package cmd

import (
	"cmd/config"
	"cmd/internal/command"
	"cmd/internal/database"
	"fmt"
	"os"
)

func Exec() {
	db := database.Setup(true)
	configuration := config.GetConfig()
	cmd := command.InitCommands(configuration, db)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
