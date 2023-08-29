package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCollectorWithEtcd() (pb.CollectorClient, func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// md := metadata.Pairs(
	// 	"hostname", system.Runtime.HostName,
	// 	"ip", system.Runtime.IP,
	// )
	// newCtx := metadata.NewOutgoingContext(ctx, md)

	var target = "etcd:///grpclb/king-collector"
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		// grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewCollectorClient(conn), func() error { return conn.Close() }, nil
}

func NewCollectorWithTarget(target string) (pb.CollectorClient, func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewCollectorClient(conn), func() error { return conn.Close() }, nil
}
