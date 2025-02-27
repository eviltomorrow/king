package model

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/feature"
)

func init() {
	domain.RegisterModel(&domain.Model{})
}

func F_01(k *chart.K) (*domain.Plan, bool) {
	if len(k.Candlesticks) == 0 {
		return nil, false
	}

	if !feature.MA150IsUP(k) {
		return nil, false
	}
	return nil, true
}
