package cmd

import (
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-ctl/cmd/metadata"
	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/infrastructure"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/resolver"
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

	RootCommand.AddCommand(metadata.ArchiveCommand)
	RootCommand.AddCommand(metadata.CrawlCommand)
}

func RunApp() error {
	RootCommand.Use = buildinfo.AppName

	component, err := infrastructure.LoadConfig(&etcd.Config{
		Endpoints:          []string{"127.0.0.1:2379"},
		ConnetTimeout:      5 * time.Second,
		StartupRetryTimes:  1,
		StartupRetryPeriod: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	if err := component.Init(); err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(component.Close)

	resolver.Register(lb.NewBuilder(etcd.Client))

	return RootCommand.Execute()
}
