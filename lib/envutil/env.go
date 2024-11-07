package envutil

import (
	"fmt"
	"path/filepath"

	"github.com/eviltomorrow/king/lib/db/mongodb"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/log"
	"github.com/eviltomorrow/king/lib/network"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/eviltomorrow/king/lib/redis"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/zlog"
)

func InitOpentrace(otel *opentrace.Config) error {
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

func InitLog(log *log.Config) error {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:  log.Level,
		Format: "json",
		File: zlog.FileLogConfig{
			Filename:    filepath.Join(system.Directory.LogDir, "data.log"),
			MaxSize:     100,
			MaxDays:     30,
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

	return nil
}

func InitNetwork(network *network.Config) error {
	system.Network.AccessIP = func() string {
		if network.AccessIP != "" {
			return network.AccessIP
		} else {
			return system.Network.BindIP
		}
	}()
	system.Network.BindIP = func() string {
		if network.BindIP != "0.0.0.0" {
			return network.BindIP
		} else {
			return system.Network.BindIP
		}
	}()

	return nil
}

func InitMySQL(c *mysql.Config) error {
	closeFunc, err := mysql.InitMySQL(c)
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(closeFunc)

	return nil
}

func InitEtcd(c *etcd.Config) error {
	closeFunc, err := etcd.InitEtcd(c)
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(closeFunc)

	return nil
}

func InitMongoDB(c *mongodb.Config) error {
	closeFunc, err := mongodb.InitMongoDB(c)
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(closeFunc)

	return nil
}

func InitRedis(c *redis.Config) error {
	closeFunc, err := redis.InitRedis(c)
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(closeFunc)

	return nil
}

func InitClientForGRPC(fs ...func() error) error {
	for _, f := range fs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
