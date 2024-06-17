package config

import (
	"path/filepath"

	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/zlog"
	jsoniter "github.com/json-iterator/go"
)

type Log struct {
	DisableTimestamp bool   `json:"disable-timestamp" toml:"disable-timestamp"`
	Level            string `json:"level" toml:"level"`
	Format           string `json:"format" toml:"format"`
	MaxSize          int    `json:"maxsize" toml:"maxsize"`
	MaxDays          int    `toml:"max-days" json:"max-days"`
	MaxBackups       int    `toml:"max-backups" json:"max-backups"`
	Compress         bool   `toml:"compress" json:"compress"`
}

func (c *Log) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

func SetupLogger(l Log) ([]func() error, error) {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:            l.Level,
		Format:           l.Format,
		DisableTimestamp: l.DisableTimestamp,
		File: zlog.FileLogConfig{
			Filename:    filepath.Join(system.Directory.VarDir, "log", "data.log"),
			MaxSize:     l.MaxSize,
			MaxDays:     l.MaxDays,
			MaxBackups:  l.MaxBackups,
			Compression: "gzip",
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
