package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/eviltomorrow/king/apps/king-web/conf"
	// "github.com/eviltomorrow/king/apps/king-web/domain/controller"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/envutil"
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/flagsutil"
	"github.com/eviltomorrow/king/lib/fs"
	"github.com/eviltomorrow/king/lib/http/server"
	"github.com/eviltomorrow/king/lib/pprofutil"
	"github.com/eviltomorrow/king/lib/procutil"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/zlog"
	flags "github.com/jessevdk/go-flags"
	"go.uber.org/zap"
)

func RunApp() error {
	_, err := flags.Parse(flagsutil.Opts)
	if err != nil {
		return err
	}

	if flagsutil.Opts.Version {
		fmt.Println(buildinfo.Version())
		os.Exit(0)
	}

	if flagsutil.Opts.Daemon {
		if err := procutil.RunAppInBackground(os.Args); err != nil {
			log.Fatalf("[F] Daemon: run app in background failure, nest error: %v", err)
		}
		return nil
	}

	if flagsutil.Opts.EnablePprof {
		go func() {
			if err := pprofutil.Run(flagsutil.Opts.PprofAddr); err != nil {
				log.Fatalf("[F] Run pprof failure, nest error: %v", err)
			}
		}()
	}
	defer func() {
		finalizer.RunCleanupFuncs()
	}()

	c, err := conf.ReadConfig(flagsutil.Opts)
	if err != nil {
		return fmt.Errorf("read config failure, nest error: %v", err)
	}

	if err := envutil.InitOpentrace(c.Otel); err != nil {
		return fmt.Errorf("init opentrace failure, nest error: %v", err)
	}
	if err := envutil.InitLog(c.Log); err != nil {
		return fmt.Errorf("init log failure, nest error: %v", err)
	}
	if err := envutil.InitNetwork(c.HTTP); err != nil {
		return fmt.Errorf("init network failure, nest error: %v", err)
	}

	s := server.NewHTTP(
		c.HTTP,
		c.Log,
		// controller.NewCron().Service(),
	)
	if err := s.Serve(); err != nil {
		return fmt.Errorf("cron serve failure, nest error: %v", err)
	}
	finalizer.RegisterCleanupFuncs(s.Stop)

	releaseFile, err := procutil.CreatePidFile()
	if err != nil {
		return fmt.Errorf("create pid file failure, nest error ;%v", err)
	}
	finalizer.RegisterCleanupFuncs(releaseFile)

	// 必须在最后
	if err := fs.RewriteStderrToFile(); err != nil {
		return fmt.Errorf("rewrite stderr to filre failure, nest error: %v", err)
	}

	zlog.Info("System info", zap.String("system", system.String()))
	zlog.Info("Config info", zap.String("config", c.String()))
	zlog.Info("App start success", zap.String("version", buildinfo.MainVersion), zap.String("commited-id", buildinfo.GitSha))

	procutil.StopDaemon()
	procutil.WaitForSigterm()
	zlog.Info("App stop complete", zap.String("launched-time", system.LaunchTime()))
	return nil
}