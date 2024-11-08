package controller

import (
	"context"
	"fmt"
	"time"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

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
		MarketIndexChange: &pb.MarketIndexChange{
			ShangZheng:  status.MarketIndexChange.ShangZheng,
			ShenZheng:   status.MarketIndexChange.ShenZheng,
			ChuangYe:    status.MarketIndexChange.Chuangye,
			BeiZheng_50: status.MarketIndexChange.BeiZheng50,
			KeChuang_50: status.MarketIndexChange.KeChuang50,
		},
		MarketStockChange: &pb.MarketStockChange{
			Rise_0_1:   status.MarketStockChange.Rise_0_1,
			Rise_1_3:   status.MarketStockChange.Rise_1_3,
			Rise_3_5:   status.MarketStockChange.Rise_3_5,
			Rise_5_10:  status.MarketStockChange.Rise_5_10,
			Rise_10_15: status.MarketStockChange.Rise_10_15,
			Rise_15_30: status.MarketStockChange.Rise_15_30,
			Rise_30N:   status.MarketStockChange.Rise_30_N,

			Fell_0_1:   status.MarketStockChange.Fell_0_1,
			Fell_1_3:   status.MarketStockChange.Fell_1_3,
			Fell_3_5:   status.MarketStockChange.Fell_3_5,
			Fell_5_10:  status.MarketStockChange.Fell_5_10,
			Fell_10_15: status.MarketStockChange.Fell_10_15,
			Fell_15_30: status.MarketStockChange.Fell_15_30,
			Fell_30N:   status.MarketStockChange.Fell_30_N,
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

func (c *Finder) ReportWeek(ctx context.Context, req *wrapperspb.StringValue) (*pb.MarketStatus, error) {
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

	_ = t
	// result := make(chan *data.Stock, 64)
	// pipe := make(chan *data.Stock, 64)
	// _ = result
	// var wg sync.WaitGroup

	// var count atomic.Int64
	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		for stock := range pipe {
	// 			quotes, err := data.GetQuote(context.Background(), t, stock.Code, data.DAY)
	// 			if err != nil {
	// 				zlog.Error("GetQuote failure", zap.Error(err), zap.String("code", stock.Code))
	// 				continue
	// 			} else {
	// 				k, err := chart.NewK(context.Background(), stock, quotes)
	// 				if err != nil {
	// 					zlog.Error("New k failure", zap.Error(err), zap.String("code", stock.Code))
	// 					continue
	// 				}

	// 				score, ok, err := service.ScanModel(k)
	// 				if err != nil {
	// 					zlog.Error("ScanModel failure", zap.Error(err), zap.String("code", stock.Code))
	// 					continue
	// 				}

	// 				if ok {
	// 					count.Add(1)
	// 					fmt.Printf("info: %d, %s\r\n", score, stock.Name)
	// 				}
	// 			}
	// 		}

	// 		wg.Done()
	// 	}()
	// }

	// go func() {
	// 	if err := data.FetchStock(context.Background(), pipe); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// wg.Wait()

	// fmt.Println(count.Load())
	return nil, nil
}
