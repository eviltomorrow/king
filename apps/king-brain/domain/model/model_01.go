package model

import (
	"fmt"
	"log"

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

	days := []int{200}
	k.CalMaMany(days)

	next, err := chart.CalculateMaOnNext(k, 200, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(next)

	return &domain.Plan{K: k}, nil
}

func forecastDirectionOnNext(k *chart.K, ma []float64, day int, count int) ([]string, error) {
	var (
		currentMA     = k.Candlesticks[len(k.Candlesticks)-1].Indicators.Trend.Ma[day]
		currentClosed = k.Candlesticks[len(k.Candlesticks)-1].Close
	)
	return nil, nil
}
