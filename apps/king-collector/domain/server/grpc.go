package server

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-collector/domain/service"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
	"github.com/eviltomorrow/king/lib/grpc/server"
	"github.com/eviltomorrow/king/lib/opentrace"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPC struct {
	EtcdClient *clientv3.Client
	Host       string
	Port       int
	AppName    string

	bootstrap *server.Bootstrap

	pb.UnimplementedCollectorServer
}

func (g *GRPC) CrawlMetadata(ctx context.Context, req *wrapperspb.StringValue) (*pb.CrawlCounter, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, source is nil")
	}
	if req.Value != "sina" && req.Value != "net126" {
		return nil, fmt.Errorf("invalid request, source is %s", req.Value)
	}

	ctx, span := opentrace.DefaultTracer().Start(ctx, "DataQuick")
	defer span.End()

	span.SetAttributes(attribute.String("req", req.Value))
	total, ignore, err := service.SynchronizeMetadataQuick(ctx, req.Value)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return &pb.CrawlCounter{Total: total, Ignore: ignore}, nil
}

func (g *GRPC) StoreMetadata(ctx context.Context, req *wrapperspb.StringValue) (*pb.StoreCounter, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, date is nil")
	}

	ctx, span := opentrace.DefaultTracer().Start(ctx, "StoreMetadataToStorage")
	defer span.End()

	total, stock, day, week, err := service.PushMetadataToStorage(ctx, req.Value)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return &pb.StoreCounter{Total: total, Stock: stock, Day: day, Week: week}, nil
}

func (g *GRPC) Startup() error {
	g.bootstrap = server.NewGrpcBootstrap(
		server.WithListenHost(g.Host),
		server.WithPort(g.Port),
		server.WithAppName(g.AppName),
		server.WithEtcdClient(g.EtcdClient),
		server.WithRegisterServerFunc(func(s *grpc.Server) {
			pb.RegisterCollectorServer(s, g)
		}),
	)

	return g.bootstrap.Init()
}

func (g *GRPC) Stop() error {
	if g.bootstrap != nil {
		return g.bootstrap.Stop()
	}
	return nil
}
