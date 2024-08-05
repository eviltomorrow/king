package tests

import (
	"log"

	"github.com/eviltomorrow/king/apps/king-brain/conf"
	"github.com/eviltomorrow/king/lib/envutil"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/flagsutil"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/infrastructure"
	"google.golang.org/grpc/resolver"
)

func init() {
	c := conf.InitializeDefaultConfig(&flagsutil.Flags{DisableStdlog: false})

	if err := envutil.InitBaseComponent(c.Otel, c.Log, c.GRPC); err != nil {
		log.Fatal(err)
	}

	for _, ic := range []infrastructure.Config{c.Etcd} {
		component, err := infrastructure.LoadConfig(ic)
		if err != nil {
			log.Fatalf("load config failure, nest error: %v, name: %s", err, ic.Name())
		}

		if err := component.Init(); err != nil {
			log.Fatalf("component init failure, nest error: %v", err)
		}
		finalizer.RegisterCleanupFuncs(component.Close)
	}

	resolver.Register(lb.NewBuilder(etcd.Client))
}
