package client

import (
	"fmt"

	"github.com/eviltomorrow/king/lib/finalizer"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

var DefaultFinder pb.FinderClient

func InitBrain() error {
	client, shutdown, err := NewBrainWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefaultFinder = client
	return nil
}

func NewBrainWithEtcd() (pb.FinderClient, func() error, error) {
	target := "etcd:///grpclb/king-brain"
	conn, err := grpc.NewClient(
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewFinderClient(conn), func() error { return conn.Close() }, nil
}

func NewBrainWithTarget(target string) (pb.FinderClient, func() error, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewFinderClient(conn), func() error { return conn.Close() }, nil
}
