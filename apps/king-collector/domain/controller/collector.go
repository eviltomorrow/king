package controller

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-collector/conf"
	"github.com/eviltomorrow/king/apps/king-collector/domain/metadata"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
	"github.com/eviltomorrow/king/lib/opentrace"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Collector struct {
	pb.UnimplementedCollectorServer

	config *conf.Collector
}

func NewCollector(config *conf.Collector) *Collector {
	return &Collector{config: config}
}

func (c *Collector) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterCollectorServer(server, c)
	}
}

func (c *Collector) CrawlMetadata(ctx context.Context, req *wrapperspb.StringValue) (*pb.CrawlMetadataResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, source is nil")
	}
	if req.Value != "sina" && req.Value != "net126" {
		return nil, fmt.Errorf("invalid request, source is %s", req.Value)
	}

	ctx, span := opentrace.DefaultTracer().Start(ctx, "DataQuick")
	defer span.End()

	span.SetAttributes(attribute.String("req", req.Value))
	total, ignore, err := metadata.SynchronizeMetadataQuick(ctx, req.Value, c.config.CodeList, c.config.RandomPeriod)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return &pb.CrawlMetadataResponse{Total: total, Ignore: ignore}, nil
}

func (c *Collector) StoreMetadata(ctx context.Context, req *wrapperspb.StringValue) (*pb.StoreMetadataResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, date is nil")
	}

	ctx, span := opentrace.DefaultTracer().Start(ctx, "StoreMetadataToStorage")
	defer span.End()

	stock, quote, err := metadata.StoreMetadataToStorage(ctx, req.Value)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return &pb.StoreMetadataResponse{Affected: &pb.StoreMetadataResponse_AffectedCount{
		Stock: stock,
		Quote: quote,
	}}, nil
}
