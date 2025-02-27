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
	if len(k.Candlesticks) == 0 {
		return false
	}

	last := k.Candlesticks[len(k.Candlesticks)-1]
	ma150 := last.Indicators.Trend.MA[chart.DAY150]

	if last.High > ma150 && last.Low < ma150 {
		return true
	}
	return false
}
