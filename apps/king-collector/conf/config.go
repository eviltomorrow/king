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
	MongoDB   *config.MongoDB   `json:"mongodb" toml:"mongodb" mapstructure:"mongodb"`
	Etcd      *config.Etcd      `json:"etcd" toml:"etcd" mapstructure:"etcd"`
	Log       *config.Log       `json:"log" toml:"log" mapstructure:"log"`
	GRPC      *config.GRPC      `json:"grpc" toml:"grpc" mapstructure:"grpc"`
	Otel      *config.Opentrace `json:"otel" toml:"otel" mapstructure:"otel"`
	Collector *Collector        `json:"collector" toml:"collector" mapstructure:"collector"`
}

type Collector struct {
	CodeList     []string `json:"code_list" toml:"code_list" mapstructure:"code_list"`
	Source       string   `json:"source" toml:"source" mapstructure:"source"`
	Crontab      string   `json:"crontab" toml:"crontab" mapstructure:"crontab"`
	RandomPeriod []int    `json:"random_period" toml:"random_period" mapstructure:"random_period"`
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
	Etcd: &config.Etcd{
		Endpoints: []string{
			"127.0.0.1:2379",
		},
		ConnetTimeout:      5 * time.Second,
		StartupRetryTimes:  3,
		StartupRetryPeriod: 5 * time.Second,
	},
	MongoDB: &config.MongoDB{
		DSN: "mongodb://admin:admin123@mongo:27017/transaction_db",

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
	Log: &config.Log{
		Level:         "info",
		DisableStdlog: false,
	},
	GRPC: &config.GRPC{
		AccessIP:   "",
		BindIP:     "0.0.0.0",
		BindPort:   50003,
		DisableTLS: true,
	},
	Collector: &Collector{
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
		RandomPeriod: []int{5, 30},
	},
}
