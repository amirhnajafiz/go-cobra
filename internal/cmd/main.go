package cmd

import (
	"github.com/amirhnajafiz/go-cobra/internal/cmd/crypto"
	"github.com/amirhnajafiz/go-cobra/internal/cmd/migration"
	"github.com/amirhnajafiz/go-cobra/internal/cmd/server"
	"github.com/spf13/cobra"
)

// Execute will setup database, configurations and cobra
func Execute() {
	rootCmd := cobra.Command{}

	rootCmd.AddCommand(
		crypto.GetCommand(),
		migration.GetCommand(),
		server.GetCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
