package model

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func init() {
	domain.RegisterModel(&domain.Model{Desc: "", F: F_01})
}

func F_01(k *chart.K) (*domain.Plan, error) {
	if len(k.Candlesticks) <= 200 {
		return nil, nil
	}

	days := []int{10, 20, 50, 150, 200}
	k.CalMaMany(days)

	var (
		cache = make([]*chart.Candlestick, 0, len(k.Candlesticks))
		next  = make([]*chart.Candlestick, 0, 3)
	)
	for _, c := range k.Candlesticks {
		cache = append(cache, c)
		if len(cache) == 0 {
			continue
		}
		_ = next
	}
	return &domain.Plan{K: k}, nil
}

var _ = forecastDirectionOnNext

func forecastDirectionOnNext(k *chart.K, next []float64, day int) ([]string, error) {
	if len(k.Candlesticks) < day+1 {
		return nil, fmt.Errorf("no enough candlesticks")
	}

	current := k.Candlesticks[len(k.Candlesticks)-1]

	ma, ok := current.Indicators.Trend.Ma[day]
	if !ok {
		return nil, fmt.Errorf("not found ma")
	}

	fmt.Printf("day: %d, %f, %f, %v\r\n", day, current.Closed, ma, next)
	return nil, nil
}
