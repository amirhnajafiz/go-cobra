package cmd

import (
	"cmd/config"
	"cmd/internal/database"
	"cmd/server/handler"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
)

func serverCmd(configuration config.Config, db *gorm.DB) *cobra.Command {
	return &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle Subsequent requests
			handler.HandleRequests(configuration, db)
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

	cmd.SetUsageTemplate(`[33mUsage:[0m{{if .Runnable}}
		  	{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
		  	{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
		
			Aliases:
			  {{.NameAndAliases}}{{end}}{{if .HasExample}}
			
			Examples:
			{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
			
			[33mAvailable Commands:[0m {{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
			  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
			
			[33mFlags:[0m 
			{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}
			
			Global Flags:
			{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
			
			Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
			  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
			
			Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
	`)

	cmd.AddCommand(serverCmd(configuration, db))

	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
