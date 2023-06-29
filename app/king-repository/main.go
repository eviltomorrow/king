package main

import (
	"log"
	"os"

	"github.com/eviltomorrow/king/app/king-repository/cmd"
	"github.com/eviltomorrow/king/app/king-repository/cmd/metadata"
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
	Short: "Stock trade data repository",
	Long:  "King-repository is a stock trade data repository",
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

func runApp() error {
	return RootCommand.Execute()
}

func main() {
	initCommand()

	if err := runApp(); err != nil {
		log.Printf("[F] Run app command(embeded) failure, nest error: %v", err)
		os.Exit(1)
	}
}
