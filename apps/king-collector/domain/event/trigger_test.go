package event

import (
	"context"
	"log"
	"testing"

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
	// defer client.Close()

	resolver.Register(lb.NewBuilder(client))
}

func TestNotifyWithEmail(t *testing.T) {
	if err := SendEmail(context.Background(), "test", "no matter"); err != nil {
		t.Fatalf("Notify failure, nest error: %v", err)
	}
}
