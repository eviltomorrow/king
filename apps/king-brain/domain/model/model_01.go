package model

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/feature"
)

func init() {
	domain.RegisterModel(&domain.Model{Desc: "", F: F_01})
}

func F_01(k *chart.K) (*domain.Plan, bool) {
	if len(k.Candlesticks) <= 200 {
		return nil, false
	}

	k.CalMoreMA([]int{150, 200})

	for _, f := range []func(k *chart.K) bool{
		feature.MA150IsUP,
		feature.MA150NearByClosed,
		feature.MA150GtMA200,
	} {
		if ok := f(k); !ok {
			return nil, false
		}
	}

	return &domain.Plan{K: k}, true
}
