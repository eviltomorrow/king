package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/conf"
	"github.com/eviltomorrow/king/apps/king-collector/domain/event"
	"github.com/eviltomorrow/king/apps/king-collector/domain/server"
	"github.com/eviltomorrow/king/apps/king-collector/domain/service"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/cleanup"
	"github.com/eviltomorrow/king/lib/config"
	"github.com/eviltomorrow/king/lib/db/mongodb"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/fs"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/grpc/middleware"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/eviltomorrow/king/lib/procutil"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/workflow"
	"github.com/eviltomorrow/king/lib/zlog"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

var workflowsFunc = func() []workflow.Job {
	workflow.Register("setRuntimeEnv", setRuntimeEnv)
	workflow.Register("loadConfig", loadConfig)
	workflow.Register("printCfg", printCfg)
	workflow.Register("setGlobalVars", setGlobalVars)
	workflow.Register("runOpentrace", runOpentrace)
	workflow.Register("runDB", runDB)
	workflow.Register("runServer", runServer)
	workflow.Register("runCron", runCron)
	workflow.Register("buildPidFile", buildPidFile)
	workflow.Register("rewritePaniclog", rewritePaniclog)
	workflow.Register("notifyStopDaemon", notifyStopDaemon)
	return workflow.Finish()
}()

var cfg = conf.Default
var (
	isBackground bool
	ppid         int
)

func init() {
	StartCommand.PersistentFlags().BoolVarP(&isBackground, "daemon", "d", false, "run app in background")
	StartCommand.PersistentFlags().IntVar(&ppid, "ppid", -1, "bootup parent pid")
}

var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start the app in teminal(-d in background)",
	Run: func(cmd *cobra.Command, args []string) {
		begin := time.Now()
		if isBackground {
			if err := procutil.RunAppBackground(os.Args[0], []string{"start"}); err != nil {
				log.Fatalf("[F] Run app in background failure, nest error: %v", err)
			}
			return
		}

		defer func() {
			for _, err := range cleanup.RunCleanupFuncs() {
				zlog.Error("Run cleanup funcs failure", zap.Error(err))
			}
			zlog.Info("Stop app complete", zap.String("app-name", buildinfo.AppName), zap.String("running-duration", system.Runtime.RunningDuration()))
		}()

		for _, job := range workflowsFunc {
			if err := job.F(); err != nil {
				log.Fatalf("[F] Run job failure, nest error: %v, job name: %v", err, job.Name)
			}
		}

		zlog.Info("Start app success", zap.String("app-name", buildinfo.AppName), zap.Duration("cost", time.Since(begin)))
		procutil.WaitForSigterm()
	},
}

func setRuntimeEnv() error {
	for _, dir := range []string{
		filepath.Join(system.Runtime.RootDir, "/var/log"),
		filepath.Join(system.Runtime.RootDir, "/var/run"),
	} {
		if err := fs.CreateDir(dir); err != nil {
			return fmt.Errorf("create dir failure, nest error: %v", err)
		}
	}
	return nil
}

func loadConfig() error {
	if err := cfg.LoadFile(filepath.Join(system.Runtime.RootDir, "/etc/config.toml")); err != nil {
		return err
	}

	closeFuncs, err := config.SetupLogger(cfg.Log)
	if err != nil {
		return err
	}
	cleanup.RegisterCleanupFuncs(closeFuncs...)

	return nil
}

func setGlobalVars() error {
	middleware.LogDir = filepath.Join(system.Runtime.RootDir, "/var/log")
	etcd.Endpoints = cfg.Etcd.Endpoints
	mongodb.DSN = cfg.MongoDB.DSN
	event.Source = cfg.Collector.Source
	service.CodeList = cfg.Collector.CodeList
	if strings.Count(cfg.Collector.RandomPeriod, ",") == 1 {
		attrs := strings.Split(cfg.Collector.RandomPeriod, ",")
		v1, _ := strconv.Atoi(attrs[0])
		v2, _ := strconv.Atoi(attrs[1])
		if v1 < v2 && v1 > 0 && v2 < 100 {
			service.RandomPeriod = [2]int{v1, v2}
		}
	}
	opentrace.OtelDSN = cfg.Otel.DSN
	return nil
}

func runOpentrace() error {
	if !cfg.Otel.Enable {
		return nil
	}
	shutdown, err := opentrace.InitTraceProvider()
	if err != nil {
		return fmt.Errorf("init trace provider failure, nest error: %v", err)
	}
	cleanup.RegisterCleanupFuncs(shutdown)
	return nil
}

func runDB() error {
	if err := mongodb.Connect(); err != nil {
		return err
	}
	cleanup.RegisterCleanupFuncs(mongodb.Close)

	return nil
}

func runServer() error {
	client, err := etcd.NewClient()
	if err != nil {
		return err
	}
	cleanup.RegisterCleanupFuncs(client.Close)

	if err := middleware.InitLogger(); err != nil {
		return err
	}
	resolver.Register(lb.NewBuilder(client))

	g := &server.GRPC{
		AppName:    buildinfo.AppName,
		Host:       cfg.Server.Host,
		Port:       cfg.Server.Port,
		EtcdClient: client,
	}

	if err := g.Startup(); err != nil {
		return err
	}
	cleanup.RegisterCleanupFuncs(g.Stop)

	return nil
}

func runCron() error {
	c := cron.New()
	_, err := c.AddFunc(cfg.Collector.Crontab, func() {
		ctx, span := opentrace.DefaultTracer().Start(context.Background(), "RunCron")
		defer span.End()

		if err := event.TriggerFetchMetadataEveryWeekDay(ctx); err != nil {
			span.SetStatus(codes.Error, "FetchMetadataEveryWeekDay failure")
			span.RecordError(err)
			zlog.Error("Crontab run archive metadata every weekday failure", zap.Error(err))
		}
	})
	if err != nil {
		return err
	}
	c.Start()

	cleanup.RegisterCleanupFuncs(func() error {
		c.Stop()
		return nil
	})
	return nil
}

func rewritePaniclog() error {
	fs.StderrFilePath = filepath.Join(system.Runtime.RootDir, "/var/log/panic.log")
	if err := fs.RewriteStderrFile(); err != nil {
		zlog.Error("RewriteStderrFile failure", zap.Error(err))
	}
	return nil
}

func buildPidFile() error {
	closeFunc, err := procutil.CreatePidFile(filepath.Join(system.Runtime.RootDir, fmt.Sprintf("/var/run/%s.pid", buildinfo.AppName)), system.Runtime.Pid)
	if err != nil {
		return err
	}
	cleanup.RegisterCleanupFuncs(closeFunc)
	return nil
}

func notifyStopDaemon() error {
	if ppid == -1 {
		return nil
	}

	if _, err := procutil.FindProcessWithPid(ppid); err != nil {
		return err
	}

	confirmationBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading confirmation bytes from stdin: %v", err)
	}
	if len(confirmationBytes) == 0 {
		return nil
	}

	return procutil.NotifyStopDaemon(confirmationBytes)
}

func printCfg() error {
	zlog.Info("Load config success", zap.String("config", cfg.String()))
	return nil
}
