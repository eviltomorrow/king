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
	if len(k.Candlesticks) == 0 {
		return nil, false
	}

	if !feature.MA150IsUP(k) {
		return nil, false
	}

	if !feature.MA150NearByClosed(k) {
		return nil, false
	}

	if !feature.MA150GtMA200(k) {
		return nil, false
	}

	return &domain.Plan{K: k}, true
}
