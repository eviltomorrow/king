package cmd

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-ctl/cmd/metadata"
	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Short: "king-ctl tool for king-*(all service).",
	Long:  "King 程序辅助终端",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	RootCommand.SetHelpCommand(&cobra.Command{
		Use:    "",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf(`unknown command "help" for "%s"`, cmd.Root().Name())
		},
	})

	RootCommand.AddCommand(metadata.CrawlCommand)
	RootCommand.AddCommand(metadata.StatsCommand)
	RootCommand.AddCommand(metadata.StoreCommand)
}

func RunApp() error {
	RootCommand.Use = buildinfo.AppName

	return RootCommand.Execute()
}
