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

	segment, err := chart.CalculateMaToSegment(k, 150)
	if err != nil {
		return nil, err
	}

	fmt.Println(segment)
	return &domain.Plan{K: k}, nil
}
