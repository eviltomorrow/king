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

	days := []int{20, 30, 40, 50, 100, 150, 200}
	k.CalMaMany(days)

	for _, day := range days {
		data, err := chart.CalculateMaOnNext(k, day, 10)
		if err == chart.ErrNoData {
			continue
		}
		if err != nil {
			return nil, err
		}
		fmt.Println(data)
	}

	return &domain.Plan{K: k}, nil
}
