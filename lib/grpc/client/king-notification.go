package client

import (
	"fmt"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewEmailWithEtcd() (pb.EmailClient, func() error, error) {
	target := "etcd:///grpclb/king-notification"
	conn, err := grpc.NewClient(
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewEmailClient(conn), func() error { return conn.Close() }, nil
}

func NewEmailWithTarget(target string) (pb.EmailClient, func() error, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewEmailClient(conn), func() error { return conn.Close() }, nil
}

func NewNtfyWithEtcd() (pb.NtfyClient, func() error, error) {
	target := "etcd:///grpclb/king-notification"
	conn, err := grpc.NewClient(
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewNtfyClient(conn), func() error { return conn.Close() }, nil
}

func NewNtfyWithTarget(target string) (pb.NtfyClient, func() error, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewNtfyClient(conn), func() error { return conn.Close() }, nil
}
