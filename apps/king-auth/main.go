package main

import (
	"log"

	"github.com/eviltomorrow/king/apps/king-auth/cmd"
	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/system"
)

var (
	AppName     = "unknown"
	MainVersion = "unknown"
	GitSha      = "unknown"
	BuildTime   = "unknown"
)

func init() {
	buildinfo.AppName = AppName
	buildinfo.MainVersion = MainVersion
	buildinfo.GitSha = GitSha
	buildinfo.BuildTime = BuildTime
}

func main() {
	if err := system.LoadRuntime(); err != nil {
		log.Fatalf("[F] App: load system runtime failure, nest error: %v", err)
	}

	if err := cmd.RunApp(); err != nil {
		log.Fatalf("[F] App: run app failure, nest error: %v", err)
	}
}
