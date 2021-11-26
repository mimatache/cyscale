package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mimatache/cyscale/internal/info"
)

var Version = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	RunE: func(cmd *cobra.Command, args []string) error {
		appInfo := info.AppInfo()
		fmt.Println("Name:      ", appInfo.Name)
		fmt.Println("Version:   ", appInfo.Version)
		fmt.Println("Hash:      ", appInfo.Hash)
		fmt.Println("Build Date:", appInfo.BuildDate)
		return nil
	},
}
