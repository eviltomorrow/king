package tests

import (
	"time"

	"github.com/eviltomorrow/king/apps/king-auth/conf"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/flagsutil"
	"github.com/eviltomorrow/king/lib/log"
	"github.com/eviltomorrow/king/lib/network"
	"github.com/eviltomorrow/king/lib/opentrace"
)

func InitializeDefaultConfig(opts *flagsutil.Flags) *conf.Config {
	return &conf.Config{
		Etcd: &etcd.Config{
			Endpoints: []string{
				"192.168.33.10:2379",
			},
			ConnectTimeout:     5 * time.Second,
			StartupRetryTimes:  3,
			StartupRetryPeriod: 5 * time.Second,
		},
		Otel: &opentrace.Config{
			Enable:         false,
			DSN:            "192.168.33.10:4317",
			ConnectTimeout: 5 * time.Second,
		},
		Log: &log.Config{
			Level:         "info",
			DisableStdlog: opts.DisableStdlog,
		},
		GRPC: &network.Config{
			AccessIP:   "",
			BindIP:     "0.0.0.0",
			BindPort:   50005,
			DisableTLS: true,
		},
	}
}
