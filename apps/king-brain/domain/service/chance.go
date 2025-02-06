package service

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func FindPossibleChance(k *chart.K) {
	domain.ScanFeatureFunc(k)
}
