package config

import jsoniter "github.com/json-iterator/go"

type Etcd struct {
	Endpoints []string `json:"endpoints" toml:"endpoints"`
}

func (c *Etcd) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}
