package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/eviltomorrow/king/app/king-collector/conf"
	"github.com/eviltomorrow/king/app/king-collector/server"
	"github.com/eviltomorrow/king/app/king-collector/service"
	"github.com/eviltomorrow/king/app/king-collector/service/synchronize"
	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/cleanup"
	"github.com/eviltomorrow/king/lib/db/mongodb"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/fs"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/grpc/middleware"
	"github.com/eviltomorrow/king/lib/procutil"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/zlog"
	"github.com/robfig/cron/v3"
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
	runCron,
	runServer,
	buildPidFile,
	rewritePaniclog,
	stopDaemonStarter,
}

var cfg = conf.Default
var (
	isDaemonVar bool
	ppidVar     int
)

func init() {
	StartCommand.PersistentFlags().BoolVarP(&isDaemonVar, "daemon", "d", false, "run app in background")
	StartCommand.PersistentFlags().IntVar(&ppidVar, "ppid", -1, "the pid of parent progress")
	StartCommand.Flags().MarkHidden("ppid")
}

var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start the app in teminal(-d in background)",
	Run: func(cmd *cobra.Command, args []string) {
		var begin = time.Now()
		if isDaemonVar {
			if err := procutil.BackgroundStartup(buildinfo.AppName, []string{"start", "--ppid", fmt.Sprintf("%d", os.Getpid())}, nil, nil); err != nil {
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

func setRuntimeEnv() error {
	for _, dir := range []string{
		filepath.Join(system.Runtime.RootDir, "/log"),
		filepath.Join(system.Runtime.RootDir, "/var/run"),
	} {
		if err := fs.CreateDir(dir); err != nil {
			return fmt.Errorf("create dir failure, nest error: %v", err)
		}
	}
	return nil
}

func loadConfig() error {
	if err := cfg.LoadFile(filepath.Join(system.Runtime.RootDir, "/etc/global.toml")); err != nil {
		return err
	}

	closeFuncs, err := conf.SetupLogger(cfg.Log)
	if err != nil {
		return err
	}
	cleanup.RegisterCleanupFuncs(closeFuncs...)

	return nil
}

func setGlobalVars() error {
	middleware.LogDir = filepath.Join(system.Runtime.RootDir, "/log")
	etcd.Endpoints = cfg.Etcd.Endpoints
	mongodb.DSN = cfg.MongoDB.DSN
	service.Source = cfg.Collector.Source
	synchronize.CodeList = cfg.Collector.CodeList
	if strings.Count(cfg.Collector.RandomPeriod, ",") == 1 {
		var attrs = strings.Split(cfg.Collector.RandomPeriod, ",")
		v1, _ := strconv.Atoi(attrs[0])
		v2, _ := strconv.Atoi(attrs[1])
		if v1 < v2 && v1 > 0 && v2 < 100 {
			synchronize.RandomPeriod = [2]int{v1, v2}
		}
	}
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

	zlog.Info("Start server(grpc) success", zap.String("host", g.Host), zap.Int("port", g.Port))
	return nil
}

func runCron() error {
	var c = cron.New()
	_, err := c.AddFunc(cfg.Collector.Crontab, func() {
		if err := service.ArchiveMetadataEveryWeekDay(); err != nil {
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
	fs.StderrFilePath = filepath.Join(system.Runtime.RootDir, "/log/panic.log")
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

func stopDaemonStarter() error {
	if ppidVar == -1 || ppidVar == 0 || ppidVar == 1 || ppidVar == 2 {
		return nil
	}
	if os.Getppid() != ppidVar {
		return fmt.Errorf("panic: the app not started with -d")
	}
	process, err := os.FindProcess(ppidVar)
	if err != nil {
		return err
	}
	if err := process.Signal(syscall.SIGUSR1); err != nil {
		return fmt.Errorf("panic: notify process with signal SIGUSER1 failure, nest error: %v", err)
	}
	return nil
}

func printCfg() error {
	zlog.Info("Load config success", zap.String("config", cfg.String()))
	return nil
}
