package envutil

import (
	"fmt"
	"path/filepath"

	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/grpc/middleware"
	"github.com/eviltomorrow/king/lib/grpc/server"
	"github.com/eviltomorrow/king/lib/log"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/zlog"
)

func InitBaseComponent(otel *opentrace.Config, log *log.Config, server *server.Config) error {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:  log.Level,
		Format: "json",
		File: zlog.FileLogConfig{
			Filename:    filepath.Join(system.Directory.LogDir, "data.log"),
			MaxSize:     100,
			MaxDays:     90,
			MaxBackups:  90,
			Compression: "gzip",
		},
		DisableStacktrace: true,
		DisableStdlog:     log.DisableStdlog,
	})
	if err != nil {
		return fmt.Errorf("init global log failure, nest error: %v", err)
	}
	zlog.ReplaceGlobals(global, prop)
	finalizer.RegisterCleanupFuncs(global.Sync)

	// middleware log
	midlog, err := middleware.InitLogger(&zlog.Config{
		Level:  log.Level,
		Format: "json",
		File: zlog.FileLogConfig{
			Filename:    filepath.Join(system.Directory.LogDir, "access.log"),
			MaxSize:     100,
			MaxDays:     90,
			MaxBackups:  90,
			Compression: "gzip",
		},
		DisableStacktrace: true,
		DisableStdlog:     log.DisableStdlog,
	})
	if err != nil {
		return fmt.Errorf("init midlog failure, nest error: %v", err)
	}
	finalizer.RegisterCleanupFuncs(midlog)

	system.Network.AccessIP = func() string {
		if server.AccessIP != "" {
			return server.AccessIP
		} else {
			return system.Network.BindIP
		}
	}()
	system.Network.BindIP = func() string {
		if server.BindIP != "0.0.0.0" {
			return server.BindIP
		} else {
			return system.Network.BindIP
		}
	}()

	shutdown, err := opentrace.InitTraceProvider(&opentrace.Config{
		DSN:            otel.DSN,
		ConnectTimeout: otel.ConnectTimeout,
	})
	if err != nil {
		return fmt.Errorf("init opentrace failure, nest error: %v", err)
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	return nil
}
