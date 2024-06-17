package conf

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/eviltomorrow/king/lib/config"
	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	Redis Redis `json:"redis" toml:"redis"`
	MySQL MySQL `json:"mysql" toml:"mysql"`

	Etcd   config.Etcd      `json:"etcd" toml:"etcd"`
	Log    config.Log       `json:"log" toml:"log"`
	Server config.Server    `json:"server" toml:"server"`
	Otel   config.Opentrace `json:"otel" toml:"otel"`

	Global Global `json:"global" toml:"global"`
}

type Global struct {
	AccessTokenExpiresIn  time.Duration `json:"access-token-expires-in" toml:"access-token-expires-in"`
	RefreshTokenExpiresIn time.Duration `json:"refresh-token-expires-in" toml:"refresh-token-expires-in"`

	TokenLimitPerAccount int64 `json:"token-limit-per-account" toml:"token-limit-per-account"`
}

type Redis struct {
	DSN string `json:"dsn" toml:"dsn"`
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

var Default = Config{
	Etcd: config.Etcd{
		Endpoints: []string{
			"127.0.0.1:2379",
		},
	},
	Redis: Redis{
		DSN: "redis://:admin123@redis:6379/0?protocol=3",
	},
	MySQL: MySQL{
		DSN:     "admin:admin123@tcp(127.0.0.1:3306)/king_auth?charset=utf8mb4&parseTime=true&loc=Local",
		MinOpen: 3,
		MaxOpen: 10,
	},
	Otel: config.Opentrace{
		Enable: true,
		DSN:    "otel-collector:4317",
	},
	Log: config.Log{
		DisableTimestamp: false,
		Level:            "info",
		Format:           "json",
		MaxSize:          100,
		MaxDays:          90,
		MaxBackups:       90,
		Compress:         true,
	},
	Server: config.Server{
		Host: "0.0.0.0",
		Port: 5277,
	},
	Global: Global{
		AccessTokenExpiresIn:  10 * time.Minute,
		RefreshTokenExpiresIn: 120 * time.Minute,

		TokenLimitPerAccount: 10,
	},
}
