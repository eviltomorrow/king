package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-email"
	"github.com/eviltomorrow/king/lib/system"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func NewEmailWithEtcd() (pb.EmailClient, func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	md := metadata.Pairs(
		"hostname", system.Runtime.HostName,
		"ip", system.Runtime.IP,
	)
	newCtx := metadata.NewOutgoingContext(ctx, md)

	var target = "etcd:///grpclb/king-email"
	conn, err := grpc.DialContext(
		newCtx,
		target,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewEmailClient(conn), func() error { return conn.Close() }, nil
}

func NewEmailWithTarget(target string) (pb.EmailClient, func() error, error) {
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
	return pb.NewEmailClient(conn), func() error { return conn.Close() }, nil
}
