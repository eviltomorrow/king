package conf

import (
	"time"

	"github.com/eviltomorrow/king/lib/config"
	"github.com/eviltomorrow/king/lib/flagsutil"
	"github.com/eviltomorrow/king/lib/log"
	"github.com/eviltomorrow/king/lib/network"
	"github.com/eviltomorrow/king/lib/opentrace"
	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	HTTP *network.Config   `json:"http" toml:"http" mapstructure:"http"`
	Log  *log.Config       `json:"log" toml:"log" mapstructure:"log"`
	Otel *opentrace.Config `json:"otel" toml:"otel" mapstructure:"otel"`
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
		c.Otel.VerifyConfig,
		c.Log.VerifyConfig,
		c.HTTP.VerifyConfig,
	} {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func InitializeDefaultConfig(opts *flagsutil.Flags) *Config {
	return &Config{
		Otel: &opentrace.Config{
			Enable:         true,
			DSN:            "127.0.0.1:4317",
			ConnectTimeout: 5 * time.Second,
		},
		Log: &log.Config{
			Level:         "info",
			DisableStdlog: opts.DisableStdlog,
		},
		HTTP: &network.Config{
			AccessIP:   "",
			BindIP:     "0.0.0.0",
			BindPort:   50005,
			DisableTLS: true,
		},
	}
}
