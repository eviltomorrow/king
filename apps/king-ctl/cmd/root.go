package cmd

import (
	"time"

	"github.com/eviltomorrow/king/apps/king-ctl/cmd/metadata"
	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/config"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/infrastructure"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/resolver"
)

var RootCommand = &cobra.Command{
	Short: "king-ctl tool for king-*(all service).",
	Long:  "king-ctl tool for king-*(all service), support for king-collector and so on.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	RootCommand.AddCommand(metadata.ArchiveCommand)
	RootCommand.AddCommand(metadata.CrawlCommand)
}

func RunApp() error {
	RootCommand.Use = buildinfo.AppName

	component, err := infrastructure.LoadConfig(&config.Etcd{
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
	resolver.Register(lb.NewBuilder(etcd.Client))

	return RootCommand.Execute()
}
