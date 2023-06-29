package config

type Etcd struct {
	Endpoints []string `json:"endpoints" toml:"endpoints"`
}
