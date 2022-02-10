package command

import (
	"cmd/internal/cmd/server"
	"cmd/internal/config"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func ServerCmd(configuration config.Config, db *gorm.DB) *cobra.Command {
	set := server.Setup{
		Configuration: configuration,
		DB:            db,
	}
	return &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle Subsequent requests
			set.HandleRequests()
		},
	}
}
