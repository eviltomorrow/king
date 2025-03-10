package model

import (
	"log"

	"github.com/eviltomorrow/king/lib/envutil"
	"github.com/eviltomorrow/king/lib/flagsutil"
	"github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/tests"
)

func init() {
	c := tests.InitializeDefaultConfig(&flagsutil.Flags{DisableStdlog: false})

	if err := envutil.InitLog(c.Log); err != nil {
		log.Fatalf("init log failure, nest error: %v", err)
	}
	if err := envutil.InitNetwork(c.GRPC); err != nil {
		log.Fatalf("init network failure, nest error: %v", err)
	}
	if err := envutil.InitEtcd(c.Etcd); err != nil {
		log.Fatalf("init etcd failure, nest error: %v", err)
	}
	if err := envutil.InitClientForGRPCWithTarget("192.168.33.10:50001", client.InitStorageWithTarget); err != nil {
		log.Fatalf("init etcd failure, nest error: %v", err)
	}
}
