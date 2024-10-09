package client

import (
	"fmt"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-cron"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewSchedulerWithEtcd() (pb.SchedulerClient, func() error, error) {
	target := "etcd:///grpclb/king-cron"
	conn, err := grpc.NewClient(
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewSchedulerClient(conn), func() error { return conn.Close() }, nil
}

func NewSchedulerWithTarget(target string) (pb.SchedulerClient, func() error, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewSchedulerClient(conn), func() error { return conn.Close() }, nil
}
