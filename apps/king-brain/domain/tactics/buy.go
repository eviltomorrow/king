package tactics

import "github.com/eviltomorrow/king/apps/king-brain/domain/chart"

func TryToFindBuyPoint(k *chart.K) (float64, bool) {
	if len(k.Candlesticks) == 0 {
		return 0, false
	}

	// last := k.Candlesticks[len(k.Candlesticks)-1]

	return 0, false
}
