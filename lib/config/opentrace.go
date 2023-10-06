package config

import jsoniter "github.com/json-iterator/go"

type Opentrace struct {
	Enable bool   `json:"enable" toml:"enable"`
	DSN    string `json:"dsn" toml:"dsn"`
}

func (c *Opentrace) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}
