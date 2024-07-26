package config

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Opentrace struct {
	Enable        bool          `json:"enable" toml:"enable" mapstructure:"enable"`
	DSN           string        `json:"dsn" toml:"dsn" mapstructure:"dsn"`
	ConnetTimeout time.Duration `json:"connect_timeout" toml:"-" mapstructure:"-"`
}

func (c *Opentrace) Validate() error {
	if c.DSN == "" {
		return fmt.Errorf("otel.dsn has no value")
	}
	if c.ConnetTimeout <= 0 {
		return fmt.Errorf("otel.connect_timeout has no value")
	}
	return nil
}

func (c *Opentrace) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}
