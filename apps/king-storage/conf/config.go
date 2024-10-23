package conf

import (
	"time"

	"github.com/eviltomorrow/king/lib/config"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/flagsutil"
	"github.com/eviltomorrow/king/lib/log"
	"github.com/eviltomorrow/king/lib/network"
	"github.com/eviltomorrow/king/lib/opentrace"
	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	Etcd  *etcd.Config      `json:"etcd" toml:"etcd" mapstructure:"etcd"`
	Log   *log.Config       `json:"log" toml:"log" mapstructure:"log"`
	MySQL *mysql.Config     `json:"mysql" toml:"mysql" mapstructure:"mysql"`
	GRPC  *network.Config   `json:"grpc" toml:"grpc" mapstructure:"grpc"`
	Otel  *opentrace.Config `json:"otel" toml:"otel" mapstructure:"otel"`
}

func (c *Config) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

func ReadConfig(opts *flagsutil.Flags) (*Config, error) {
	c := InitializeDefaultConfig(opts)

	if err := config.ReadFile(c, opts.ConfigFile); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) IsConfigValid() error {
	for _, f := range []func() error{
		c.Etcd.VerifyConfig,
		c.MySQL.VerifyConfig,
		c.Otel.VerifyConfig,
		c.Log.VerifyConfig,
		c.GRPC.VerifyConfig,
	} {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func InitializeDefaultConfig(opts *flagsutil.Flags) *Config {
	return &Config{
		Etcd: &etcd.Config{
			Endpoints: []string{
				"127.0.0.1:2379",
			},
			ConnectTimeout:     5 * time.Second,
			StartupRetryTimes:  3,
			StartupRetryPeriod: 5 * time.Second,
		},
		MySQL: &mysql.Config{
			DSN:     "admin:admin123@tcp(127.0.0.1:3306)/king_storage?charset=utf8mb4&parseTime=true&loc=Local",
			MinOpen: 3,
			MaxOpen: 10,

			MaxLifetime:        5 * time.Minute,
			ConnectTimeout:     5 * time.Second,
			StartupRetryTimes:  3,
			StartupRetryPeriod: 5 * time.Second,
		},
		Otel: &opentrace.Config{
			Enable:         true,
			DSN:            "127.0.0.1:4317",
			ConnectTimeout: 5 * time.Second,
		},
		Log: &log.Config{
			Level:         "info",
			DisableStdlog: opts.DisableStdlog,
		},
		GRPC: &network.Config{
			AccessIP:   "",
			BindIP:     "0.0.0.0",
			BindPort:   50001,
			DisableTLS: true,
		},
	}
}
