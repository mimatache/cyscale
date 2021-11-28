package about

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mimatache/cyscale/internal/info"
)

// Version gives dinformation about the current application (version, build date, hash) to be able to more easily track down the build
var Version = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		appInfo := info.AppInfo()
		fmt.Println("Name:      ", appInfo.Name)
		fmt.Println("Version:   ", appInfo.Version)
		fmt.Println("Hash:      ", appInfo.Hash)
		fmt.Println("Build Date:", appInfo.BuildDate)
		return nil
	},
}
