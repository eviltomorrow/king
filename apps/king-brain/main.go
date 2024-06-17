package main

import (
	"log"

	"github.com/eviltomorrow/king/apps/king-brain/cmd"
	"github.com/eviltomorrow/king/lib/buildinfo"
	libcmd "github.com/eviltomorrow/king/lib/cmd"
	"github.com/eviltomorrow/king/lib/system"
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
	Short: "Possible chance finder with brain",
	Long:  "king-brain is a possible chance finder with brain",
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
	if err := system.InitRuntime(); err != nil {
		log.Fatalf("[F] App: init system runtime failure, nest error: %v", err)
	}
	initCommand()

	if err := runApp(); err != nil {
		log.Fatalf("[F] App: Run app command(embeded) failure, nest error: %v", err)
	}
}
