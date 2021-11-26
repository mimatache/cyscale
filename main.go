package main

import (
	"fmt"
	"os"

	"github.com/mimatache/cyscale/internal/commands"
	"github.com/mimatache/cyscale/internal/info"
)

func main() {
	if err := commands.Root(info.AppInfo().Name).Execute(); err != nil {
		err = fmt.Errorf("could not run command; %w", err)
		fmt.Println(err)
		os.Exit(1)
	}
}
