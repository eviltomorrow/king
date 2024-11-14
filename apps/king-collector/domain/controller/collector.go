package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/conf"
	"github.com/eviltomorrow/king/apps/king-collector/domain/db"
	"github.com/eviltomorrow/king/apps/king-collector/domain/service"
	"github.com/eviltomorrow/king/lib/db/mongodb"
	"github.com/eviltomorrow/king/lib/grpc/pb/entity"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
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

func (c *Collector) CrawlMetadataAsync(ctx context.Context, req *wrapperspb.StringValue) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, scheduler_id is nil")
	}
	if req.Value == "" {
		return nil, fmt.Errorf("invalid request, scheduler_id is %s", req.Value)
	}

	go func(id string) {
		begin := time.Now()
		zlog.Info("crawl metadata begin")

		var (
			total, ignore int64
			err           error
		)
		if c.config.CrawlMode == "slow" {
			total, ignore, err = service.CrawlMetadataSlow(context.Background(), c.config.Source, c.config.CodeList, c.config.RandomPeriod)
		} else {
			total, ignore, err = service.CrawlMetadataQuick(context.Background(), c.config.Source, c.config.CodeList)
		}

		if err != nil {
			zlog.Error("crawl metadata failure", zap.Error(err))
		} else {
			zlog.Info("crawl metadata completed", zap.Int64("total", total), zap.Int64("ignore", ignore), zap.Duration("cost", time.Since(begin)))
		}

		schedulerId, err := service.Callback(id, err)
		if err != nil {
			zlog.Error("callback failure", zap.Error(err))
		} else {
			zlog.Info("callback success", zap.String("schedulerId", schedulerId))
		}
	}(req.Value)

	return &emptypb.Empty{}, nil
}

func (c *Collector) CrawlMetadata(ctx context.Context, req *wrapperspb.StringValue) (*pb.CrawlMetadataResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, source is nil")
	}
	if req.Value != "sina" {
		return nil, fmt.Errorf("invalid request, source is %s", req.Value)
	}

	total, ignore, err := service.CrawlMetadataQuick(ctx, req.Value, c.config.CodeList)
	if err != nil {
		return nil, err
	}
	return &pb.CrawlMetadataResponse{Total: total, Ignore: ignore}, nil
}

func (c *Collector) FetchMetadata(req *wrapperspb.StringValue, resp grpc.ServerStreamingServer[entity.Metadata]) error {
	if req == nil {
		return fmt.Errorf("invalid request, date is nil")
	}

	var (
		offset, limit int64 = 0, 30
		lastID        string
	)

	for {
		metadata, err := db.SelectMetadataRange(context.Background(), mongodb.DB, offset, limit, req.Value, lastID)
		if err != nil {
			return err
		}

		for _, md := range metadata {
			if err := resp.Send(&entity.Metadata{
				Source:          md.Source,
				Code:            md.Code,
				Name:            md.Name,
				Open:            md.Open,
				YesterdayClosed: md.YesterdayClosed,
				Latest:          md.Latest,
				High:            md.High,
				Low:             md.Low,
				Volume:          md.Volume,
				Account:         md.Account,
				Date:            md.Date,
				Time:            md.Time,
				Suspend:         md.Suspend,
			}); err != nil {
				return err
			}
		}
		if len(metadata) < int(limit) {
			break
		}
		offset += limit
	}

	return nil
}
