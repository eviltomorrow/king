package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/eviltomorrow/king/apps/king-cron/conf"
	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/controller"
	_ "github.com/eviltomorrow/king/apps/king-cron/domain/plan"
	"github.com/eviltomorrow/king/apps/king-cron/domain/service"
	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/envutil"
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/flagsutil"
	"github.com/eviltomorrow/king/lib/fs"
	"github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/lib/grpc/server"
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
	if err := envutil.InitNetwork(c.GRPC); err != nil {
		return fmt.Errorf("init network failure, nest error: %v", err)
	}
	if err := envutil.InitEtcd(c.Etcd); err != nil {
		return fmt.Errorf("init etcd failure, nest error: %v", err)
	}
	if err := envutil.InitMySQL(c.MySQL); err != nil {
		return fmt.Errorf("init mysql failure, nest error: %v", err)
	}

	initClientFunc := []func() error{
		client.InitBrain, client.InitCollector, client.InitStorage, client.InitEmail, client.InitNTFY, client.InitTemplate,
	}
	if err := envutil.InitClientForGRPC(initClientFunc...); err != nil {
		return fmt.Errorf("init grpc client failure, nest error: %v", err)
	}

	s := server.NewGRPC(
		c.GRPC,
		c.Log,
		controller.NewCron().Service(),
	)
	if err := s.Serve(); err != nil {
		return fmt.Errorf("cron serve failure, nest error: %v", err)
	}
	finalizer.RegisterCleanupFuncs(s.Stop)

	cron := service.NewScheduler()
	for _, c := range c.Crons {
		plan, ok := domain.GetPlan(c.Plan)
		zlog.Debug("Plan will be register", zap.String("name", plan.GetAlias()))
		if ok {
			if err := cron.Register(c.Crontab, plan); err != nil {
				return err
			}
		}
	}
	if err := cron.Start(); err != nil {
		return fmt.Errorf("cron start failure, nest error: %v", err)
	}

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
