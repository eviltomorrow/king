package controller

import (
	"context"
	"fmt"
	"time"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	_ "github.com/eviltomorrow/king/apps/king-brain/domain/model"
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

func ReportDaily(ctx context.Context, req *wrapperspb.StringValue) (*pb.StatsInfo, error) {
	return nil, nil
}

func ReportWeek(ctx context.Context, req *wrapperspb.StringValue) (*pb.StatsInfo, error) {
	return nil, nil
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
