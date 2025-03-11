package model

import (
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

	// days := []int{150}
	// span := 150
	// k.CalMaMany(days)

	// for _, day := range days {
	// 	data, err := chart.CalculateMaOnNext(k, day, span, 1)
	// 	if err == chart.ErrNoData {
	// 		continue
	// 	}
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return &domain.Plan{K: k}, nil
}
