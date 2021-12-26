package command

import (
	"cmd/config"
	"cmd/server/handler"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func ServerCmd(configuration config.Config, db *gorm.DB) *cobra.Command {
	return &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle Subsequent requests
			handler.HandleRequests(configuration, db)
		},
	}
}
