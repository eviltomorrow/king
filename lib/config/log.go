package config

type Log struct {
	DisableTimestamp bool   `json:"disable-timestamp" toml:"disable-timestamp"`
	Level            string `json:"level" toml:"level"`
	Format           string `json:"format" toml:"format"`
	MaxSize          int    `json:"maxsize" toml:"maxsize"`
	MaxDays          int    `toml:"max-days" json:"max-days"`
	MaxBackups       int    `toml:"max-backups" json:"max-backups"`
	Dir              string `toml:"dir" json:"dir"`
	Compress         bool   `toml:"compress" json:"compress"`
}
