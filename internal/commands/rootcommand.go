package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mimatache/cyscale/internal/commands/version"
)

// Root is the start of the application. All other commands should either be added to this command, or to a command that is in intself added to this command
func Root(app string) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   app,
		Short: fmt.Sprintf("%s is used to find basic security misconfigurations, given the provided input.", app),
	}

	rootCommand.AddCommand(
		version.Version,
	)

	return rootCommand
}
