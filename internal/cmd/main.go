package cmd

import (
	"cmd/internal/database"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func serverCmd() *cobra.Command {
	return &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle Subsequent requests
			// handleRequests()
			return nil
		},
	}
}

func Run() {
	cmd := &cobra.Command{
		Use:     "dispatch",
		Short:   "Dispatch Server",
		Version: "0.1",
	}

	cmd.AddCommand(serverCmd())

	db := database.Setup(true)
	log.Println(db.Error.Error())

	if err := cmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(0)
	}
}
