package cmd

import (
	"cmd/config"
	"cmd/internal/database"
	"cmd/server/handler"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
)

func serverCmd(configuration config.Config, db *gorm.DB) *cobra.Command {
	return &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle Subsequent requests
			handler.HandleRequests(configuration, db)
			return nil
		},
	}
}

func Exec() {
	db := database.Setup(true)
	configuration := config.GetConfig()
	cmd := &cobra.Command{
		Use:     "dispatch",
		Short:   "Dispatch Server",
		Version: "0.1",
	}

	cmd.AddCommand(serverCmd(configuration, db))

	if err := cmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(0)
	}
}
