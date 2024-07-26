package config

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Redis struct {
	DSN                string        `json:"dsn" toml:"dsn" mapstructure:"dsn"`
	ConnetTimeout      time.Duration `json:"connect_timeout" toml:"-" mapstructure:"-"`
	StartupRetryTimes  int           `json:"startup_retry_times" toml:"-" mapstructure:"-"`
	StartupRetryPeriod time.Duration `json:"startup_retry_period" toml:"-" mapstructure:"-"`
}

func (c *Redis) Name() string {
	return "redis"
}

func (c *Redis) Validate() error {
	if c.DSN == "" {
		return fmt.Errorf("redis.dsn has no value")
	}
	if c.ConnetTimeout <= 0 {
		return fmt.Errorf("redis.connect_timeout has no value")
	}
	if c.StartupRetryTimes <= 0 {
		return fmt.Errorf("redis.startup_retry_times has no value")
	}
	if c.StartupRetryPeriod <= 0 {
		return fmt.Errorf("redis.startup_retry_period has no value")
	}
	return nil
}

func (r *Redis) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(r)
	return string(buf)
}

func (r *Redis) MarshalConfig() ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(r)
}
