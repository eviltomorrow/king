package data

import (
	"log"

	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"google.golang.org/grpc/resolver"
)

func init() {
	etcd.Endpoints = []string{"etcd:2379"}
	client, err := etcd.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	resolver.Register(lb.NewBuilder(client))
}
