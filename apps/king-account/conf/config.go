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
	GRPC   *config.GRPC      `json:"grpc" toml:"grpc" mapstructure:"grpc"`
	Otel   *config.Opentrace `json:"otel" toml:"otel" mapstructure:"otel"`
}

type Global struct {
	TransferFeeRatio     string   `json:"transfer_fee_ratio" toml:"transfer_fee_ratio" mapstructure:"transfer_fee_ratio"`
	TransferFeePayMethod []string `json:"transfer_fee_pay_method" toml:"transfer_fee_pay_method" mapstructure:"transfer_fee_pay_method"`

	CommissionRatio     string   `json:"commission_ratio" toml:"commission_ratio" mapstructure:"commission_ratio"`
	CommissionPayMethod []string `json:"commission_pay_method" toml:"commission_pay_method" mapstructure:"commission_pay_method"`

	StampTaxRatio     string   `json:"stamp_tax_ratio" toml:"stamp_tax_ratio" mapstructure:"stamp_tax_ratio"`
	StampTaxPayMethod []string `json:"stamp_tax_pay_method" toml:"stamp_tax_pay_method" mapstructure:"stamp_tax_pay_method"`
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
		TransferFeeRatio:     "0.001%",
		TransferFeePayMethod: []string{"buy", "sell"},
		CommissionRatio:      "0.018%",
		CommissionPayMethod:  []string{"buy", "sell"},
		StampTaxRatio:        "0.05%",
		StampTaxPayMethod:    []string{"sell"},
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
	Log: &config.Log{
		Level:         "info",
		DisableStdlog: false,
	},
	GRPC: &config.GRPC{
		AccessIP:   "",
		BindIP:     "0.0.0.0",
		BindPort:   50005,
		DisableTLS: true,
	},
}
