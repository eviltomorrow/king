package cmd

import (
	"fmt"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/spf13/cobra"
)

var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the app version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(buildinfo.Version())
	},
}
