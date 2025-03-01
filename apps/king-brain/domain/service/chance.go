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

func FindPossibleChance(ctx context.Context, date time.Time) []*domain.Plan {
	var (
		wg sync.WaitGroup

		pipeCh = make(chan *data.Stock, 64)
		planCh = make(chan *domain.Plan, 16)
	)

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			for stock := range pipeCh {
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

				for _, plan := range domain.ScanModel(k) {
					planCh <- plan
				}
			}

			wg.Done()
		}()
	}

	go func() {
		defer func() {
			close(planCh)
		}()

		if err := data.FetchStock(context.Background(), pipeCh); err != nil {
			zlog.Error("FetchStock failure", zap.Error(err))
		}
		wg.Wait()
	}()

	plans := make([]*domain.Plan, 0, 64)
	for plan := range planCh {
		fmt.Println(plan.K.Code, plan.K.Name)
		plans = append(plans, plan)
	}

	return plans
}
