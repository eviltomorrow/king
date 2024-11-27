package controller

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	_ "github.com/eviltomorrow/king/apps/king-brain/domain/model"
	"github.com/eviltomorrow/king/apps/king-brain/domain/service"
)

type Finder struct {
	pb.UnimplementedFinderServer
}

func NewFinder() *Finder {
	return &Finder{}
}

func (c *Finder) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterFinderServer(server, c)
	}
}

func (c *Finder) reportMarketStatus(ctx context.Context, date time.Time, kind string) (*pb.MarketStatus, error) {
	status, err := service.ReportMarketStatus(ctx, date, kind)
	if err != nil {
		return nil, err
	}

	result := &pb.MarketStatus{
		Date: status.Date,
		Week: status.Week,
		MarketIndex: &pb.MarketIndex{
			ShangZheng: &pb.Point{
				Value:      status.MarketIndex.ShangZheng.Value,
				HasChanged: status.MarketIndex.ShangZheng.HasChanged,
			},
			ShenZheng: &pb.Point{
				Value:      status.MarketIndex.ShenZheng.Value,
				HasChanged: status.MarketIndex.ShenZheng.HasChanged,
			},
			ChuangYe: &pb.Point{
				Value:      status.MarketIndex.ChuangYe.Value,
				HasChanged: status.MarketIndex.ChuangYe.HasChanged,
			},
			BeiZheng_50: &pb.Point{
				Value:      status.MarketIndex.BeiZheng50.Value,
				HasChanged: status.MarketIndex.BeiZheng50.HasChanged,
			},
			KeChuang_50: &pb.Point{
				Value:      status.MarketIndex.KeChuang50.Value,
				HasChanged: status.MarketIndex.KeChuang50.HasChanged,
			},
		},
		MarketStockCount: &pb.MarketStockCount{
			Total:     status.MarketStockCount.Total,
			Rise:      status.MarketStockCount.Rise,
			RiseGt_7:  status.MarketStockCount.RiseGT7,
			RiseGt_15: status.MarketStockCount.RiseGT15,
			Fell:      status.MarketStockCount.Fell,
			FellGt_7:  status.MarketStockCount.FellGT7,
			FellGt_15: status.MarketStockCount.FellGT15,
		},
	}
	return result, nil
}

func (c *Finder) ReportDaily(ctx context.Context, req *wrapperspb.StringValue) (*pb.MarketStatus, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, req is nil")
	}

	t, err := time.ParseInLocation(time.DateOnly, req.Value, time.Local)
	if err != nil {
		return nil, fmt.Errorf("parse time failure, nest error: %v", err)
	}

	return c.reportMarketStatus(ctx, t, "day")
}

func (c *Finder) ReportWeekly(ctx context.Context, req *wrapperspb.StringValue) (*pb.MarketStatus, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, req is nil")
	}

	t, err := time.ParseInLocation(time.DateOnly, req.Value, time.Local)
	if err != nil {
		return nil, fmt.Errorf("parse time failure, nest error: %v", err)
	}

	return c.reportMarketStatus(ctx, t, "week")
}

func (c *Finder) FindPossibleChance(ctx context.Context, req *wrapperspb.StringValue) (*pb.Chances, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, req is nil")
	}

	t, err := time.ParseInLocation(time.DateOnly, req.Value, time.Local)
	if err != nil {
		return nil, fmt.Errorf("parse time failure, nest error: %v", err)
	}

	result := make(chan *data.Stock, 64)
	pipe := make(chan *data.Stock, 64)
	_ = result
	var wg sync.WaitGroup

	var count atomic.Int64
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for stock := range pipe {
				quotes, err := data.GetQuotesN(context.Background(), t, stock.Code, "day", 250)
				if err != nil {
					zlog.Error("GetQuote failure", zap.Error(err), zap.String("code", stock.Code))
					continue
				} else {
					k, err := chart.NewK(context.Background(), stock, quotes)
					if err != nil {
						zlog.Error("New k failure", zap.Error(err), zap.String("code", stock.Code))
						continue
					}

					score, ok, err := service.ScanModel(k)
					if err != nil {
						zlog.Error("ScanModel failure", zap.Error(err), zap.String("code", stock.Code))
						continue
					}

					if ok {
						count.Add(1)
						fmt.Printf("info: %d, %s\r\n", score, stock.Name)
					}
				}
			}

			wg.Done()
		}()
	}

	go func() {
		if err := data.FetchStock(context.Background(), pipe); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	fmt.Println(count.Load())
	return nil, nil
}
