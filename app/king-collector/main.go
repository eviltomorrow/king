package main

import (
	"log"

	"github.com/eviltomorrow/king/app/king-collector/cmd"
	"github.com/eviltomorrow/king/app/king-collector/cmd/metadata"
	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/spf13/cobra"
)

var (
	AppName     = "unknown"
	MainVersion = "unknown"
	GitSha      = "unknown"
	BuildTime   = "unknown"
)

var RootCommand = &cobra.Command{
	Use:   AppName,
	Short: "Crawl stock data(every weekday)",
	Long:  "King-collector crawl stock data from [sina/net126] every weekday",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	buildinfo.AppName = AppName
	buildinfo.MainVersion = MainVersion
	buildinfo.GitSha = GitSha
	buildinfo.BuildTime = BuildTime
}

func initCommand() {
	RootCommand.AddCommand(metadata.RootCommand)
	RootCommand.AddCommand(cmd.StartCommand)
	RootCommand.AddCommand(cmd.StopCommand)
	RootCommand.AddCommand(cmd.VersionCommand)
	RootCommand.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
}

func main() {
	initCommand()

	if err := runApp(); err != nil {
		log.Fatalf("[F] Run app command(embeded) failure, nest error: %v", err)
	}
}

func runApp() error {
	return RootCommand.Execute()
}
