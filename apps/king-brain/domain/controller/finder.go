package controller

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	"github.com/eviltomorrow/king/apps/king-brain/domain/service"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
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

func (c *Finder) DiscoverPossibleChance(ctx context.Context, req *wrapperspb.StringValue) (*pb.Chances, error) {
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
				quotes, err := data.GetQuote(context.Background(), t, stock.Code, data.DAY)
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
						zlog.Info("", zap.Int("score", score), zap.String("name", stock.Name))
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
