package mongodb

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	DEFAULT_MONGODB_MIN_OPEN       = 3
	DEFAULT_MONGODB_MAX_OPEN       = 10
	DEFAULT_MONGODB_MIN_OPEN_LIMIT = DEFAULT_MONGODB_MIN_OPEN
	DEFAULT_MONGODB_MAX_OPEN_LIMIT = DEFAULT_MONGODB_MAX_OPEN
)

type Config struct {
	DSN     string `json:"dsn" toml:"dsn" mapstructure:"dsn"`
	MinOpen uint64 `json:"min_open" toml:"min_open" mapstructure:"min_open"`
	MaxOpen uint64 `json:"max_open" toml:"max_open" mapstructure:"max_open"`

	MaxLifetime        time.Duration `json:"max_lifetime" toml:"-" mapstructure:"-"`
	ConnetTimeout      time.Duration `json:"connect_timeout" toml:"-" mapstructure:"-"`
	StartupRetryTimes  int           `json:"startup_retry_times" toml:"-" mapstructure:"-"`
	StartupRetryPeriod time.Duration `json:"startup_retry_period" toml:"-" mapstructure:"-"`
}

func (c *Config) Name() string {
	return "mongodb"
}

func (c *Config) VerifyConfig() error {
	if c.DSN == "" {
		return fmt.Errorf("mongodb.dsn has no value")
	}

	if c.MinOpen == 0 {
		c.MinOpen = DEFAULT_MONGODB_MIN_OPEN
	}
	if c.MaxOpen == 0 {
		c.MaxOpen = DEFAULT_MONGODB_MAX_OPEN
	}
	if c.MinOpen < DEFAULT_MONGODB_MIN_OPEN_LIMIT {
		return fmt.Errorf("mongodb.min_open must be greater than %d", DEFAULT_MONGODB_MIN_OPEN_LIMIT)
	}
	if c.MaxOpen > DEFAULT_MONGODB_MAX_OPEN_LIMIT {
		return fmt.Errorf("mongodb.max_open must be less than %d", DEFAULT_MONGODB_MAX_OPEN_LIMIT)
	}
	if c.MinOpen > c.MaxOpen {
		return fmt.Errorf("mongodb.min_open should be less than mongodb.max-open")
	}
	if c.MaxLifetime <= 0 {
		return fmt.Errorf("mongodb.max_lifetime has no value")
	}
	if c.ConnetTimeout <= 0 {
		return fmt.Errorf("mongodb.connect_timeout has no value")
	}
	if c.StartupRetryTimes <= 0 {
		return fmt.Errorf("mongodb.startup_retry_times has no value")
	}
	if c.StartupRetryPeriod <= 0 {
		return fmt.Errorf("mongodb.startup_retry_period has no value")
	}
	return nil
}

func (m *Config) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
	return string(buf)
}

func (m *Config) MarshalConfig() ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
}
