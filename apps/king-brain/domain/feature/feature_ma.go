package feature

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func MA150IsUP(k *chart.K) bool {
	trend := chart.CalculateTrendWithMA(chart.DAY150, k)
	if len(trend) == 0 {
		return false
	}

	last := trend[len(trend)-1]
	if len(last) <= 2 {
		return false
	}

	if last[len(last)-1] > last[0] && last[len(last)-1] > last[len(last)-2] {
		return true
	}
	return false
}

func MA150NearByClosed(k *chart.K) bool {
	return false
}
