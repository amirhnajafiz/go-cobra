package command

import (
	"cmd/internal/cmd/server"
	"github.com/spf13/cobra"
)

var message = `[33mUsage:[0m{{if .Runnable}}
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

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`

func (c Commander) InitCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dispatch",
		Short:   "Dispatch Server",
		Version: "0.1",
	}
	cmd.SetUsageTemplate(message)
	cmd.AddCommand(c.ServerCmd())
	return cmd
}

func (c Commander) ServerCmd() *cobra.Command {
	set := server.Setup{
		Configuration: c.Configuration,
		DB:            c.DB,
		Logger:        c.Logger.Named("server"),
	}
	return &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle Subsequent requests
			set.HandleRequests()
		},
	}
}
