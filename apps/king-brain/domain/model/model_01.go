package model

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func init() {
	domain.RegisterModel(&domain.Model{})
}

func F_01(k *chart.K) (*domain.Plan, bool) {
	if len(k.Candlesticks) == 0 {
		return nil, false
	}

	for _, c := range k.Candlesticks {
		fmt.Println(c)
	}
	return nil, false
}
