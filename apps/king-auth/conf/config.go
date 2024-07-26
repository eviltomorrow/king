package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/eviltomorrow/king/lib/config"
	"github.com/eviltomorrow/king/lib/flagsutil"
	"github.com/eviltomorrow/king/lib/system"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

type Config struct {
	Global *Global           `json:"global" toml:"global" mapstructure:"global"`
	Etcd   *config.Etcd      `json:"etcd" toml:"etcd" mapstructure:"etcd"`
	Log    *config.Log       `json:"log" toml:"log" mapstructure:"log"`
	MySQL  *config.MySQL     `json:"mysql" toml:"mysql" mapstructure:"mysql"`
	Redis  *config.Redis     `json:"redis" toml:"redis" mapstructure:"redis"`
	GRPC   *config.GRPC      `json:"grpc" toml:"grpc" mapstructure:"grpc"`
	Otel   *config.Opentrace `json:"otel" toml:"otel" mapstructure:"otel"`
}

type Global struct {
	AccessTokenExpiresIn  time.Duration `json:"access_token_expires_in" toml:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Duration `json:"refresh_token_expires_in" toml:"refresh_token_expires_in"`

	TokenLimitPerAccount int64 `json:"token_limit_per_account" toml:"token_limit_per_account"`
}

func (c *Config) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

func ReadConfig(opts *flagsutil.Flags) (*Config, error) {
	DefaultConfig.Log.DisableStdlog = opts.DisableStdlog

	var findConfigFile = func(path string) (string, error) {
		for _, p := range []string{
			path,
			filepath.Join(system.Directory.EtcDir, "config.toml"),
		} {
			fi, err := os.Stat(p)
			if err == nil && !fi.IsDir() {
				return p, nil
			}
		}
		return "", fmt.Errorf("not found config file")
	}

	configFile, err := findConfigFile(opts.ConfigFile)
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := v.Unmarshal(DefaultConfig); err != nil {
		return nil, err
	}
	if err := DefaultConfig.validate(); err != nil {
		return nil, err
	}

	return DefaultConfig, nil
}

func (c *Config) validate() error {
	for _, f := range []func() error{
		c.Etcd.Validate,
		c.MySQL.Validate,
		c.Redis.Validate,
		c.Otel.Validate,
		c.Log.Validate,
		c.GRPC.Validate,
	} {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

var DefaultConfig = &Config{
	Global: &Global{
		AccessTokenExpiresIn:  10 * time.Minute,
		RefreshTokenExpiresIn: 120 * time.Minute,

		TokenLimitPerAccount: 10,
	},
	Etcd: &config.Etcd{
		Endpoints: []string{
			"127.0.0.1:2379",
		},
		ConnetTimeout:      5 * time.Second,
		StartupRetryTimes:  3,
		StartupRetryPeriod: 5 * time.Second,
	},
	MySQL: &config.MySQL{
		DSN:     "admin:admin123@tcp(127.0.0.1:3306)/king_storage?charset=utf8mb4&parseTime=true&loc=Local",
		MinOpen: 3,
		MaxOpen: 10,

		MaxLifetime:        5 * time.Minute,
		ConnetTimeout:      5 * time.Second,
		StartupRetryTimes:  3,
		StartupRetryPeriod: 5 * time.Second,
	},
	Otel: &config.Opentrace{
		Enable:        true,
		DSN:           "127.0.0.1:4317",
		ConnetTimeout: 5 * time.Second,
	},
	Redis: &config.Redis{
		DSN:                "redis://:admin123@127.0.0.1:6379/0?protocol=3",
		ConnetTimeout:      5 * time.Second,
		StartupRetryTimes:  3,
		StartupRetryPeriod: 5 * time.Second,
	},
	Log: &config.Log{
		Level:         "info",
		DisableStdlog: false,
	},
	GRPC: &config.GRPC{
		AccessIP:   "",
		BindIP:     "0.0.0.0",
		BindPort:   50004,
		DisableTLS: true,
	},
}
