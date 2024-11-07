package client

import (
	"fmt"

	"github.com/eviltomorrow/king/lib/finalizer"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

var DefalutCollector pb.CollectorClient

func InitCollector() error {
	client, shutdown, err := NewCollectorWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefalutCollector = client
	return nil
}

func NewCollectorWithEtcd() (pb.CollectorClient, func() error, error) {
	target := "etcd:///grpclb/king-collector"
	conn, err := grpc.NewClient(
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewCollectorClient(conn), func() error { return conn.Close() }, nil
}

func NewCollectorWithTarget(target string) (pb.CollectorClient, func() error, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewCollectorClient(conn), func() error { return conn.Close() }, nil
}
