/*
 * ------------------------------------------------------------
 * "THE BEERWARE LICENSE" (Revision 42):
 * <Mihai Matache (matache91mh@gmail.com)> wrote this code. As long as you retain this
 * notice, you can do whatever you want with this stuff. If we
 * meet someday, and you think this stuff is worth it, you can
 * buy me a beer in return.
 * ------------------------------------------------------------
 */

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
