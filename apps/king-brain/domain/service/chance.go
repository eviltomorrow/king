package service

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

type Chance struct{}

func FindPossibleChance(k *chart.K) (*Chance, bool) {
	if len(k.Candlesticks) == 0 {
		return nil, false
	}

	return nil, false
}
