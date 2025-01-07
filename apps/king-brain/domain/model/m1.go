package model

import (
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/service"
)

type ConfigM1 struct {
	SlidingWindows int
	ObserveWindows int
}

var C1 = ConfigM1{
	SlidingWindows: 10,
	ObserveWindows: 50,
}

func init() {
	service.RegisterModel(&service.Model{
		Name: "",
		Desc: "",
		F:    M1,
	})
}

func M1(k *chart.K) (*service.Strategy, bool) {
	if len(k.Candlesticks) == 0 {
		return nil, false
	}

	for i := 0; i < len(k.Candlesticks); i++ {
		c := k.Candlesticks[i]
		fmt.Printf("percent: %v, closed: %v, ma10: %v, ma50: %v, ma150: %v, date: %v\r\n", c.Volatility.PercentageChange, c.Close, c.MA[chart.Ma_10], c.MA[chart.Ma_50], c.MA[chart.Ma150], c.Date.Format(time.DateOnly))
	}

	var (
		ma_50 = make([]float64, 0, len(k.Candlesticks))
		ma150 = make([]float64, 0, len(k.Candlesticks))
	)
	for i := 0; i < len(k.Candlesticks); i++ {
		c := k.Candlesticks[i]

		v_50, ok := c.MA[chart.Ma_50]
		if ok {
			ma_50 = append(ma_50, v_50)
		}
		v150, ok := c.MA[chart.Ma150]
		if ok {
			ma150 = append(ma150, v150)
		}
	}
	return nil, false
}
