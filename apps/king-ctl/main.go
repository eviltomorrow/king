package main

import (
	"log"

	"github.com/eviltomorrow/king/apps/king-ctl/cmd"
	"github.com/eviltomorrow/king/lib/buildinfo"
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
	if err := cmd.RunApp(); err != nil {
		log.Fatalf("send reqeust to service failure, nest error: %v", err)
	}
}
