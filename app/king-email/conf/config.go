package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/eviltomorrow/king/lib/config"
	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	SmtpFile string           `json:"smtp-file" toml:"smtp-file"`
	Etcd     config.Etcd      `json:"etcd" toml:"etcd"`
	Log      config.Log       `json:"log" toml:"log"`
	Server   config.Server    `json:"server" toml:"server"`
	Otel     config.Opentrace `json:"otel" toml:"otel"`
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
		SmtpFile: "etc/smtp.json",
		Etcd: config.Etcd{
			Endpoints: []string{
				"127.0.0.1:2379",
			},
		},
		Log: config.Log{
			DisableTimestamp: false,
			Level:            "info",
			Format:           "json",
			MaxSize:          100,
			MaxDays:          180,
			MaxBackups:       90,
			Compress:         true,
		},
		Server: config.Server{
			Host: "0.0.0.0",
			Port: 5273,
		},
		Otel: config.Opentrace{
			DSN: "otel-collector:4317",
		},
	}
)
