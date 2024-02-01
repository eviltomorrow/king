package config

import jsoniter "github.com/json-iterator/go"

type MySQL struct {
	DSN     string `json:"dsn" toml:"dsn"`
	MinOpen int    `json:"min-open" toml:"min-open"`
	MaxOpen int    `json:"max-open" toml:"max-open"`
}

func (m *MySQL) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
	return string(buf)
}
