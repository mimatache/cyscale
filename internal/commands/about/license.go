package about

import (
	"fmt"

	"github.com/spf13/cobra"
)

const license = `
/*
 * ------------------------------------------------------------
 * "THE BEERWARE LICENSE" (Revision 42):
 * <Mihai Matache (matache91mh@gmail.com)> wrote this code. As long as you retain this 
 * notice, you can do whatever you want with this stuff. If we
 * meet someday, and you think this stuff is worth it, you can
 * buy me a beer in return.
 * ------------------------------------------------------------
 */
`

var License = &cobra.Command{
	Use:   "license",
	Short: "Show license information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(license)
		return nil
	},
}
