package service

import (
	"context"
	"sync"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	"github.com/eviltomorrow/king/lib/zlog"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type MarketStatus struct {
	Date string

	MarketIndexChange MarketIndexChange
	MarketStockChange MarketStockChange
}

func (m *MarketStatus) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
	return string(buf)
}

type MarketIndexChange struct {
	ShangZheng float64
	ShenZheng  float64
	Chuangye   float64
	BeiZheng50 float64
	KeChuang50 float64
}

type MarketStockChange struct {
	Rise_0_1   float64
	Rise_1_3   float64
	Rise_3_5   float64
	Rise_5_10  float64
	Rise_10_15 float64
	Rise_15_30 float64
	Rise_30_N  float64

	Fell_0_1   float64
	Fell_1_3   float64
	Fell_3_5   float64
	Fell_5_10  float64
	Fell_10_15 float64
	Fell_15_30 float64
	Fell_30_N  float64
}

func ReportMarketStatus(ctx context.Context, date time.Time, kind string) (*MarketStatus, error) {
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
	status := &MarketStatus{
		Date: t,
	}
	for r := range result {
		var lastCandlestick *chart.Candlestick
		if len(r.Candlesticks) != 0 {
			lastCandlestick = r.Candlesticks[len(r.Candlesticks)-1]
		}
		if lastCandlestick != nil && t == lastCandlestick.Date.Format(time.DateOnly) {
			switch r.Name {
			case "北证50":
				status.MarketIndexChange.BeiZheng50 = lastCandlestick.Volatility.PercentageChange
			case "科创50":
				status.MarketIndexChange.KeChuang50 = lastCandlestick.Volatility.PercentageChange
			case "上证指数":
				status.MarketIndexChange.ShangZheng = lastCandlestick.Volatility.PercentageChange
			case "深证成指":
				status.MarketIndexChange.ShenZheng = lastCandlestick.Volatility.PercentageChange
			case "创业板指":
				status.MarketIndexChange.Chuangye = lastCandlestick.Volatility.PercentageChange
			default:
				switch {
				case 0.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 1.0:
					status.MarketStockChange.Rise_0_1 += 1
				case 1.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 3.0:
					status.MarketStockChange.Rise_1_3 += 1
				case 3.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 5.0:
					status.MarketStockChange.Rise_3_5 += 1
				case 5.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 10.0:
					status.MarketStockChange.Rise_5_10 += 1
				case 10.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 15.0:
					status.MarketStockChange.Rise_10_15 += 1
				case 15.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 30.0:
					status.MarketStockChange.Rise_15_30 += 1
				case 30 <= lastCandlestick.Volatility.PercentageChange:
					status.MarketStockChange.Rise_30_N += 1

				case -1.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < 0:
					status.MarketStockChange.Fell_0_1 += 1
				case -3.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < -1.0:
					status.MarketStockChange.Fell_1_3 += 1
				case -5.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < -3.0:
					status.MarketStockChange.Fell_3_5 += 1
				case -10.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < -5.0:
					status.MarketStockChange.Fell_5_10 += 1
				case -15.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < -10.0:
					status.MarketStockChange.Fell_10_15 += 1
				case -30.0 <= lastCandlestick.Volatility.PercentageChange && lastCandlestick.Volatility.PercentageChange < -15.0:
					status.MarketStockChange.Fell_15_30 += 1
				case lastCandlestick.Volatility.PercentageChange < -30.0:
					status.MarketStockChange.Fell_30_N += 1
				default:

				}

			}
		}
	}

	return status, nil
}
