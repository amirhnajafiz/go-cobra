package cmd

import (
	"github.com/spf13/cobra"
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

	if err := cmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(0)
	}
}
