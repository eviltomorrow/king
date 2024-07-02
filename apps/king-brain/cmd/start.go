package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/conf"
	"github.com/eviltomorrow/king/apps/king-brain/domain/server"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/config"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/fs"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/grpc/middleware"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/eviltomorrow/king/lib/procutil"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/workflow"
	"github.com/eviltomorrow/king/lib/zlog"
	"github.com/spf13/cobra"
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
			if err := procutil.RunAppInBackground([]string{"start"}); err != nil {
				log.Fatalf("[F] Run app in background failure, nest error: %v", err)
			}
			return
		}

		defer func() {
			for _, err := range finalizer.RunCleanupFuncs() {
				zlog.Error("Run cleanup funcs failure", zap.Error(err))
			}
			zlog.Info("Stop app complete", zap.String("app-name", buildinfo.AppName), zap.String("launched-time", system.LaunchTime()))
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
		filepath.Join(system.Directory.VarDir, "/log"),
		filepath.Join(system.Directory.VarDir, "/run"),
	} {
		if err := fs.MkdirAll(dir); err != nil {
			return fmt.Errorf("create dir failure, nest error: %v", err)
		}
	}
	return nil
}

func loadConfig() error {
	if err := cfg.LoadFile(filepath.Join(system.Directory.EtcDir, "config.toml")); err != nil {
		return err
	}

	closeFuncs, err := config.SetupLogger(cfg.Log)
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(closeFuncs...)

	return nil
}

func setGlobalVars() error {
	middleware.LogDir = filepath.Join(system.Directory.VarDir, "/log")
	etcd.Endpoints = cfg.Etcd.Endpoints
	mysql.DSN = cfg.MySQL.DSN
	mysql.MinOpen = cfg.MySQL.MinOpen
	mysql.MaxOpen = cfg.MySQL.MaxOpen

	return nil
}

func runOpentrace() error {
	shutdown, err := opentrace.InitTraceProvider()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)
	return nil
}

func runDB() error {
	if err := mysql.Connect(); err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(mysql.Close)

	return nil
}

func runServer() error {
	client, err := etcd.NewClient()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(client.Close)

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
	finalizer.RegisterCleanupFuncs(g.Stop)

	return nil
}

func rewritePaniclog() error {
	fs.StderrFilePath = filepath.Join(system.Directory.VarDir, "/log/panic.log")
	if err := fs.RewriteStderrToFile(); err != nil {
		zlog.Error("RewriteStderrToFile failure", zap.Error(err))
	}
	return nil
}

func buildPidFile() error {
	closeFunc, err := procutil.CreatePidFile(filepath.Join(system.Directory.VarDir, fmt.Sprintf("/run/%s.pid", buildinfo.AppName)), system.Process.Pid)
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(closeFunc)
	return nil
}

func notifyStopDaemon() error {
	return procutil.StopDaemon()
}

func printCfg() error {
	zlog.Info("Load config success", zap.String("config", cfg.String()))
	return nil
}
