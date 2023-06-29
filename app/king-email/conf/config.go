package conf

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/eviltomorrow/king/lib/config"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/zlog"
	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	SmtpFile string        `json:"smtp-file" toml:"smtp-file"`
	Etcd     config.Etcd   `json:"etcd" toml:"etcd"`
	Log      config.Log    `json:"log" toml:"log"`
	Server   config.Server `json:"server" toml:"server"`
}

func (c *Config) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

func (c *Config) LoadFile(path string) error {
	if _, err := toml.DecodeFile(path, c); err != nil {
		return err
	}

	return nil
}

var (
	Default = Config{
		SmtpFile: "/etc/smtp.json",
		Etcd: config.Etcd{
			Endpoints: []string{
				"127.0.0.1:2379",
			},
		},
		Log: config.Log{
			DisableTimestamp: false,
			Level:            "info",
			Format:           "text",
			MaxSize:          30,
			MaxDays:          180,
			Dir:              "log",
			Compress:         true,
		},
		Server: config.Server{
			Host: "0.0.0.0",
			Port: 5273,
		},
	}
)

func SetupLogger(l config.Log) ([]func() error, error) {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:            l.Level,
		Format:           l.Format,
		DisableTimestamp: l.DisableTimestamp,
		File: zlog.FileLogConfig{
			Filename:   filepath.Join(system.Runtime.RootDir, l.Dir, "data.log"),
			MaxSize:    l.MaxSize,
			MaxDays:    l.MaxDays,
			MaxBackups: l.MaxBackups,
			Compress:   l.Compress,
		},
		DisableStacktrace:   true,
		DisableErrorVerbose: true,
	})
	if err != nil {
		return nil, err
	}
	zlog.ReplaceGlobals(global, prop)

	return []func() error{
		func() error { return global.Sync() },
	}, nil
}
