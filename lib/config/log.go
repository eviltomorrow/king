package config

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type Log struct {
	Level         string `json:"level" toml:"level" mapstructure:"level"`
	DisableStdlog bool   `json:"disable_stdlog" toml:"-" mapstructure:"-"`
}

func (c *Log) VerifyConfig() error {
	switch c.Level {
	case "debug", "warn", "info", "error":
	default:
		return fmt.Errorf("log.level has wrong value, level: %s", c.Level)
	}
	return nil
}

func (c *Log) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}
