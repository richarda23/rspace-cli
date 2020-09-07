package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// this variable is searched/replaced by build.sh
// if this variable is renamed, the script should be updated
var rsVersion = "0.0.18-alpha2"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of RSpace CLI",
	Long:  `All software has versions. This is RSpace's`,
	Run: func(cmd *cobra.Command, args []string) {
		messageStdErr(fmt.Sprintf("This is RSpace CLI version %s", rsVersion))
	},
}
