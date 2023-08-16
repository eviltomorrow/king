package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/eviltomorrow/king/app/king-repository/conf"
	"github.com/eviltomorrow/king/app/king-repository/server"
	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/cleanup"
	"github.com/eviltomorrow/king/lib/config"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/fs"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/grpc/middleware"
	"github.com/eviltomorrow/king/lib/procutil"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/zlog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

var workflowsFunc = []func() error{
	setRuntimeEnv,
	loadConfig,
	printCfg,
	setGlobalVars,
	runDB,
	runServer,
	buildPidFile,
	rewritePaniclog,
	notifyStopDaemon,
}

var cfg = conf.Default
var (
	isDaemonVar bool
)

func init() {
	StartCommand.PersistentFlags().BoolVarP(&isDaemonVar, "daemon", "d", false, "run app in background")
}

var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start the app in teminal(-d in background)",
	Run: func(cmd *cobra.Command, args []string) {
		var begin = time.Now()
		if isDaemonVar {
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

		for _, f := range workflowsFunc {
			if err := f(); err != nil {
				log.Fatalf("[F] Run workflow failure, nest error: %v", err)
			}
		}
		zlog.Info("Start app success", zap.String("app-name", buildinfo.AppName), zap.Duration("cost", time.Since(begin)))
		procutil.WaitForSigterm()
	},
}

func loadConfig() error {
	if err := cfg.LoadFile(filepath.Join(system.Runtime.RootDir, "/etc/global.toml")); err != nil {
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
	mysql.DSN = cfg.MySQL.DSN
	mysql.MinOpen = cfg.MySQL.MinOpen
	mysql.MaxOpen = cfg.MySQL.MaxOpen

	return nil
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

func runDB() error {
	if err := mysql.Connect(); err != nil {
		return err
	}
	cleanup.RegisterCleanupFuncs(mysql.Close)

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

	var g = &server.GRPC{
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
