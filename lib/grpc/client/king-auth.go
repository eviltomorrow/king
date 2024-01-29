package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-auth"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthWithEtcd() (pb.PassportClient, func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	target := "etcd:///grpclb/king-auth"
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewPassportClient(conn), func() error { return conn.Close() }, nil
}

func NewAuthWithTarget(target string) (pb.PassportClient, func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewPassportClient(conn), func() error { return conn.Close() }, nil
}
