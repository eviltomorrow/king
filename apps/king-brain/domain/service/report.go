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

	MarketIndex      MarketIndex
	MarketStockCount MarketStockCount
}

func (m *MarketStatus) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
	return string(buf)
}

type Point struct {
	Value      float64
	HasChanged float64
}

type MarketIndex struct {
	ShangZheng Point
	ShenZheng  Point
	ChuangYe   Point
	BeiZheng50 Point
	KeChuang50 Point
}

type MarketStockCount struct {
	Total    int64
	Rise     int64
	RiseGT7  int64
	RiseGT15 int64
	Fell     int64
	FellGT7  int64
	FellGT15 int64
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
	var total int64 = 0

	for r := range result {
		var lastCandlestick *chart.Candlestick
		if len(r.Candlesticks) != 0 {
			lastCandlestick = r.Candlesticks[len(r.Candlesticks)-1]
		}
		if lastCandlestick != nil && t == lastCandlestick.Date.Format(time.DateOnly) {
			switch r.Name {
			case "北证50":
				status.MarketIndex.BeiZheng50 = Point{
					Value:      lastCandlestick.Close,
					HasChanged: lastCandlestick.Volatility.PercentageChange,
				}
			case "科创50":
				status.MarketIndex.KeChuang50 = Point{
					Value:      lastCandlestick.Close,
					HasChanged: lastCandlestick.Volatility.PercentageChange,
				}
			case "上证指数":
				status.MarketIndex.ShangZheng = Point{
					Value:      lastCandlestick.Close,
					HasChanged: lastCandlestick.Volatility.PercentageChange,
				}
			case "深证成指":
				status.MarketIndex.ShenZheng = Point{
					Value:      lastCandlestick.Close,
					HasChanged: lastCandlestick.Volatility.PercentageChange,
				}
			case "创业板指":
				status.MarketIndex.ChuangYe = Point{
					Value:      lastCandlestick.Close,
					HasChanged: lastCandlestick.Volatility.PercentageChange,
				}
			default:
				total++
				if lastCandlestick.Volatility.PercentageChange > 0 {
					status.MarketStockCount.Rise += 1
				} else {
					status.MarketStockCount.Fell += 1
				}
				switch {
				case 7.0 <= lastCandlestick.Volatility.PercentageChange:
					status.MarketStockCount.RiseGT7 += 1
				case 15.0 <= lastCandlestick.Volatility.PercentageChange:
					status.MarketStockCount.RiseGT15 += 1
				case -7.0 >= lastCandlestick.Volatility.PercentageChange:
					status.MarketStockCount.FellGT7 += 1
				case -15.0 >= lastCandlestick.Volatility.PercentageChange:
					status.MarketStockCount.FellGT15 += 1
				default:

				}
			}
		}
	}

	status.MarketStockCount.Total = total

	return status, nil
}
