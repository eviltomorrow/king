package model

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func init() {
	domain.RegisterModel(&domain.Model{Desc: "", F: F_01})
}

func F_01(k *chart.K) (*domain.Plan, bool) {
	if len(k.Candlesticks) <= 200 {
		return nil, false
	}

	k.CalMoreMA([]int{150, 200})

	fmt.Println(k.String())

	return &domain.Plan{K: k}, true
}
