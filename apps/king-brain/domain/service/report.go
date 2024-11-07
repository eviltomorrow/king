package service

import (
	"context"
	"sync"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

type StatsInfo struct {
	Date string
	Kind string

	Desc map[string]string
}

func ReportLatest(ctx context.Context, date time.Time, kind string) (*StatsInfo, error) {
	var (
		wg sync.WaitGroup

		pipe   = make(chan *data.Stock, 64)
		result = make(chan *chart.K, 64)
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for stock := range pipe {
				quotes, err := data.GetQuotesN(ctx, date, stock.Code, kind, 2)
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

	t := date.Format(time.DateOnly)
	count := make(map[string]int)
	for r := range result {
		var lastCandlestick *chart.Candlestick
		if len(r.Candlesticks) != 0 {
			lastCandlestick = r.Candlesticks[len(r.Candlesticks)-1]
		}
		if lastCandlestick != nil && t == lastCandlestick.Date.Format(time.DateOnly) {
			switch r.Name {
			case "北证50":
			case "科创50":
			case "上证指数":
			case "深证成指":
			case "创业板指":
			default:
				var desc string
				switch {
				case 0.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 1.0:
					desc = "0.0%<=~1.0%"
				case 1.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 3.0:
					desc = "1.0%<=~3.0%"
				case 3.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 5.0:
					desc = "3.0%<=~5.0%"
				case 5.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 10.0:
					desc = "5.0%<=~10.0%"
				case 10.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 15.0:
					desc = "10.0%<=~15.0%"
				case 15.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange <= 30.0:
					desc = "15.0%<=~30.0%"
				}

				val, ok := count[desc]
				if ok {
					val = val + 1
				} else {
					val = 1
				}
				count[desc] = val
			}
		}
	}

	stats := &StatsInfo{
		Date: date.Format(time.DateOnly),
		Kind: kind,
		Desc: nil,
	}

	return stats, nil
}
