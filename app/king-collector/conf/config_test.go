package conf

import "testing"

func TestLoadFile(t *testing.T) {
	Default.LoadFile("../etc/config.toml")
	select {}
}
