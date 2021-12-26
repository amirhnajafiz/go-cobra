package command

import (
	"cmd/config"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func InitCommands(configuration config.Config, db *gorm.DB) *cobra.Command {
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

	cmd.AddCommand(ServerCmd(configuration, db))

	return cmd
}
