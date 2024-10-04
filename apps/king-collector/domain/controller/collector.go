package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/conf"
	"github.com/eviltomorrow/king/apps/king-collector/domain/metadata"
	"github.com/eviltomorrow/king/lib/grpc/callback"
	pb_collector "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Collector struct {
	pb_collector.UnimplementedCollectorServer

	config *conf.Collector
}

func NewCollector(config *conf.Collector) *Collector {
	return &Collector{config: config}
}

func (c *Collector) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb_collector.RegisterCollectorServer(server, c)
	}
}

func (c *Collector) CrawlMetadataAsync(ctx context.Context, req *wrapperspb.StringValue) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, source is nil")
	}
	if req.Value != "sina" && req.Value != "net126" {
		return nil, fmt.Errorf("invalid request, source is %s", req.Value)
	}

	go func() {
		begin := time.Now()

		total, ignore, err := metadata.SynchronizeMetadataSlow(ctx, c.config.Source, c.config.CodeList, c.config.RandomPeriod)
		if err != nil {
			zlog.Error("Crawl metadata failure", zap.Error(err))
		} else {
			zlog.Info("Crawl metadata completed", zap.Int64("total", total), zap.Int64("ignore", ignore), zap.Duration("cost", time.Since(begin)))
		}

		schedulerId, err := callback.Do(ctx, err)
		if err != nil {
			zlog.Error("Callback failure", zap.Error(err))
		} else {
			zlog.Info("Callback success", zap.String("schedulerId", schedulerId))
		}
	}()

	return &emptypb.Empty{}, nil
}

func (c *Collector) CrawlMetadata(ctx context.Context, req *wrapperspb.StringValue) (*pb_collector.CrawlMetadataResponse, error) {
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
	return &pb_collector.CrawlMetadataResponse{Total: total, Ignore: ignore}, nil
}

func (c *Collector) StoreMetadata(ctx context.Context, req *wrapperspb.StringValue) (*pb_collector.StoreMetadataResponse, error) {
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
	return &pb_collector.StoreMetadataResponse{Affected: &pb_collector.StoreMetadataResponse_AffectedCount{
		Stock: stock,
		Quote: quote,
	}}, nil
}
