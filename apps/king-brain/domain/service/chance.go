package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

func FindPossibleChance(ctx context.Context, date time.Time) {
	var (
		wg sync.WaitGroup

		pipe = make(chan *data.Stock, 64)
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for stock := range pipe {
				quotes, err := data.GetQuotesN(ctx, date, stock.Code, "day", 250)
				if err != nil {
					zlog.Error("GetQuote failure", zap.Error(err), zap.String("code", stock.Code))
					continue
				}
				k, err := chart.NewK(ctx, stock, quotes)
				if err != nil {
					zlog.Error("NewK failure", zap.Error(err), zap.String("code", stock.Code))
					continue
				}

				plans := domain.ScanModel(k)
				if len(plans) != 0 {
					fmt.Println(k.Name)
				}
			}

			wg.Done()
		}()
	}

	go func() {
		if err := data.FetchStock(context.Background(), pipe); err != nil {
			zlog.Error("FetchStock failure", zap.Error(err))
			return
		}
		wg.Wait()
	}()

	time.Sleep(3 * time.Minute)
}
