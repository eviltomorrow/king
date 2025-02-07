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

		pipe   = make(chan *data.Stock, 64)
		result = make(chan *chart.K, 64)
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
				result <- k
			}

			wg.Done()
		}()
	}

	go func() {
		defer func() {
			close(result)
		}()

		if err := data.FetchStock(context.Background(), pipe); err != nil {
			zlog.Error("FetchStock failure", zap.Error(err))
			return
		}
		wg.Wait()
	}()

	for k := range result {
		desc, level := domain.CalculateFeatureFunc(k)
		fmt.Printf("%d,%s,%s,[%s]\r\n", level, k.Code, k.Name, desc)
	}
}
