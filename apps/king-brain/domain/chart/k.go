package chart

import (
	"context"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	"github.com/eviltomorrow/king/lib/mathutil"
	jsoniter "github.com/json-iterator/go"
)

// type DayNum string

// const (
// 	DAY_10 DayNum = "10"
// 	DAY_50 DayNum = "50"
// 	DAY150 DayNum = "150"
// 	DAY200 DayNum = "200"
// )

type K struct {
	Name string
	Code string

	Candlesticks []*Candlestick
}

func (k *K) CalMa(day int) {
	closed := make([]float64, 0, len(k.Candlesticks))

	for i, c := range k.Candlesticks {
		closed = append(closed, c.Closed)
		if len(closed) >= day {
			sum := mathutil.Sum(closed[i-day+1 : i+1])
			c.Indicators.Trend.Ma[day] = mathutil.Trunc4(sum / float64(day))
		}
	}
}

func (k *K) CalMaMany(day []int) {
	closed := make([]float64, 0, len(k.Candlesticks))

	for i, c := range k.Candlesticks {
		closed = append(closed, c.Closed)

		for _, d := range day {
			if len(closed) >= d {
				sum := mathutil.Sum(closed[i-d+1 : i+1])
				c.Indicators.Trend.Ma[d] = mathutil.Trunc4(sum / float64(d))
			}
		}
	}
}

func (k *K) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(k)
	return string(buf)
}

type Candlestick struct {
	Date time.Time

	High    float64
	Low     float64
	Open    float64
	Closed  float64
	Volume  int64
	Account float64

	Indicators Indicators
}

type Indicators struct {
	Trend      *Trend
	Volatility *Volatility
}

func (c *Candlestick) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

type Trend struct {
	Ma map[int]float64
}

type Volatility struct {
	PercentPriceChange      float64 // 涨幅
	PercentVolume           float64 // 量能
	PercentAmplitude        float64 // 振幅
	AverageTransactionPrice float64 // 平均成交价
}

func NewK(ctx context.Context, stock *data.Stock, quotes []*data.Quote) (*K, error) {
	candlesticks := make([]*Candlestick, 0, len(quotes))
	// closed = make([]float64, 0, len(quotes))
	for i, quote := range quotes {
		date, err := time.Parse(time.DateOnly, quote.Date)
		if err != nil {
			return nil, err
		}

		v := &Volatility{
			PercentPriceChange: mathutil.Trunc2(float64(quote.Close-quote.YesterdayClosed) / float64(quote.YesterdayClosed) * 100),
			PercentVolume: func() float64 {
				if i != 0 {
					last := quotes[i-1]

					return mathutil.Trunc2((float64(quote.Volume-last.Volume)/float64(last.Volume) + 1) * 100)
				}
				return 0
			}(),
			PercentAmplitude: func() float64 {
				return mathutil.Trunc2((quote.High - quote.Low) / quote.YesterdayClosed * 100)
			}(),
			AverageTransactionPrice: func() float64 {
				return mathutil.Trunc2(quote.Account / float64(quote.Volume))
			}(),
		}

		c := &Candlestick{
			Date:    date,
			High:    quote.High,
			Low:     quote.Low,
			Open:    quote.Open,
			Closed:  quote.Close,
			Volume:  quote.Volume,
			Account: quote.Account,

			Indicators: Indicators{
				Trend: &Trend{
					Ma: make(map[int]float64, len(quotes)+10),
				},
				Volatility: v,
			},
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
