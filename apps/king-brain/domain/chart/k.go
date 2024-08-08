package chart

import (
	"context"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	jsoniter "github.com/json-iterator/go"
)

type MaKind string

const (
	Ma_10 MaKind = "ma_10"
	Ma_50 MaKind = "ma_50"
	Ma100 MaKind = "ma100"
	Ma150 MaKind = "ma150"
	Ma200 MaKind = "ma200"
)

type K struct {
	Name string
	Code string

	Candlesticks []*Candlestick
}

func (k *K) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(k)
	return string(buf)
}

type Candlestick struct {
	Date    time.Time
	High    float64
	Low     float64
	Open    float64
	Close   float64
	Volume  uint64
	Account float64

	MA map[MaKind]float64
}

func NewK(ctx context.Context, stock *data.Stock, quotes []*data.Quote) (*K, error) {
	var (
		candlesticks = make([]*Candlestick, 0, len(quotes))

		closed = make([]float64, 0, len(quotes))
	)
	for i, quote := range quotes {
		date, err := time.Parse(time.DateOnly, quote.Date)
		if err != nil {
			return nil, err
		}
		var c = &Candlestick{
			Date:    date,
			High:    quote.High,
			Low:     quote.Low,
			Open:    quote.Open,
			Close:   quote.Close,
			Volume:  quote.Volume,
			Account: quote.Account,
			MA:      make(map[MaKind]float64, 5),
		}

		closed = append(closed, quote.Close)
		if len(closed) >= 10 {
			c.MA[Ma_10] = calculateMa(closed[i-10+1 : i+1])
		}
		if len(closed) >= 50 {
			c.MA[Ma_50] = calculateMa(closed[i-50+1 : i-1])
		}
		if len(closed) >= 100 {
			c.MA[Ma100] = calculateMa(closed[i-100+1 : i-1])
		}
		if len(closed) >= 150 {
			c.MA[Ma150] = calculateMa(closed[i-150+1 : i-1])
		}
		if len(closed) >= 200 {
			c.MA[Ma200] = calculateMa(closed[i-200+1 : i-1])
		}
		candlesticks = append(candlesticks, c)
	}

	k := &K{
		Name:         stock.Name,
		Code:         stock.Code,
		Candlesticks: candlesticks,
	}
	return k, nil
}
