package main

import (
	"log"

	"github.com/eviltomorrow/king/apps/king-notification/cmd"
	"github.com/eviltomorrow/king/lib/buildinfo"
	libcmd "github.com/eviltomorrow/king/lib/cmd"
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
	Short: "Send email service",
	Long:  "king-notification is a email sender service",
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
	RootCommand.AddCommand(cmd.StartCommand)
	RootCommand.AddCommand(libcmd.StopCommand)
	RootCommand.AddCommand(libcmd.VersionCommand)
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
		log.Fatalf("[F] Run app command(embeded) failure, nest error: %v", err)
	}
}
