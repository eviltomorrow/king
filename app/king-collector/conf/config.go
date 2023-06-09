package conf

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eviltomorrow/king/lib/config"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/zlog"
	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	Collector Collector `json:"collector" toml:"collector"`
	MongoDB   MongoDB   `json:"mongodb" toml:"mongodb"`

	Etcd   config.Etcd   `json:"etcd" toml:"etcd"`
	Log    config.Log    `json:"log" toml:"log"`
	Server config.Server `json:"server" toml:"server"`
}

type MongoDB struct {
	DSN      string `json:"dsn" toml:"dsn"`
	Database string `json:"database"`
}

type Collector struct {
	CodeList     []string `json:"code-list" toml:"code-list"`
	Source       string   `json:"source" toml:"source"`
	Crontab      string   `json:"crontab" toml:"crontab"`
	RandomPeriod string   `json:"random-period" toml:"random-period"`
}

func (c *Config) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

func (c *Config) LoadFile(path string) error {
	if _, err := toml.DecodeFile(path, c); err != nil {
		return err
	}

	var n = strings.LastIndex(c.MongoDB.DSN, "/")
	if n == len(c.MongoDB.DSN) {
		return fmt.Errorf("panic: wrong mongodb dsn: %s, possible without db", c.MongoDB.DSN)
	}
	c.MongoDB.Database = c.MongoDB.DSN[n+1:]
	return nil
}

var (
	Default = Config{
		MongoDB: MongoDB{
			DSN:      "mongodb://127.0.0.1:27017/transaction_db",
			Database: "transaction_db",
		},
		Collector: Collector{
			CodeList: []string{
				"sh688***",
				"sh605***",
				"sh603***",
				"sh601***",
				"sh600***",
				"sz300***",
				"sz0030**",
				"sz002***",
				"sz001***",
				"sz000***",
			},
			Source:       "sina",
			Crontab:      "05 18 * * MON,TUE,WED,THU,FRI",
			RandomPeriod: "20,60",
		},
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
			Port: 5271,
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
