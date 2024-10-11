package tests

import (
	"log"

	"github.com/eviltomorrow/king/apps/king-brain/conf"
	"github.com/eviltomorrow/king/lib/envutil"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/flagsutil"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"google.golang.org/grpc/resolver"
)

func init() {
	c := conf.InitializeDefaultConfig(&flagsutil.Flags{DisableStdlog: false})

	if err := envutil.InitOpentrace(c.Otel); err != nil {
		log.Fatalf("init opentrace failure, nest error: %v", err)
	}
	if err := envutil.InitLog(c.Log); err != nil {
		log.Fatalf("init log failure, nest error: %v", err)
	}
	if err := envutil.InitNetwork(c.GRPC); err != nil {
		log.Fatalf("init network failure, nest error: %v", err)
	}
	if err := envutil.InitEtcd(c.Etcd); err != nil {
		log.Fatalf("init etcd failure, nest error: %v", err)
	}

	resolver.Register(lb.NewBuilder(etcd.Client))
}
