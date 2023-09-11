package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/eviltomorrow/king/lib/config"
	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	Etcd   config.Etcd   `json:"etcd" toml:"etcd"`
	Log    config.Log    `json:"log" toml:"log"`
	MySQL  MySQL         `json:"mysql" toml:"mysql"`
	Server config.Server `json:"server" toml:"server"`
}

type MySQL struct {
	DSN     string `json:"dsn" toml:"dsn"`
	MinOpen int    `json:"min-open" toml:"min-open"`
	MaxOpen int    `json:"max-open" toml:"max-open"`
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
		Etcd: config.Etcd{
			Endpoints: []string{
				"127.0.0.1:2379",
			},
		},
		MySQL: MySQL{
			DSN:     "admin:admin123@tcp(127.0.0.1:3306)/king_brain?charset=utf8mb4&parseTime=true&loc=Local",
			MinOpen: 3,
			MaxOpen: 10,
		},
		Log: config.Log{
			DisableTimestamp: false,
			Level:            "info",
			Format:           "text",
			MaxSize:          100,
			MaxDays:          180,
			MaxBackups:       90,
			Compress:         true,
		},
		Server: config.Server{
			Host: "0.0.0.0",
			Port: 5274,
		},
	}
)
