package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mimatache/cyscale/internal/commands/about"
	"github.com/mimatache/cyscale/internal/commands/verifier"
)

// Root is the start of the application. All other commands should either be added to this command, or to a command that is in intself added to this command
func Root(app string) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   app,
		Short: fmt.Sprintf("%s is used to find basic security misconfigurations, given the provided input.", app),
		Long: `
This is a simple application that parses json files containing information about cloud topologies and performs simple security violation scans.
This should not be treated as an exhaustive security scan of your cloud environment	
`,
	}

	rootCommand.AddCommand(
		about.Version,
		about.License,
		verifier.Verify(),
	)

	return rootCommand
}
