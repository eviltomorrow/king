package feature

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func init() {
	domain.RegisterFeatureFunc(&domain.Feature{Desc: "收盘价高于 150 日均线", F: FeatureWithClosedGTMA150})
	domain.RegisterFeatureFunc(&domain.Feature{Desc: "收盘价高于 200 日均线", F: FeatureWithClosedGTMA150})
	domain.RegisterFeatureFunc(&domain.Feature{Desc: "150 日均线高于 200 日均线", F: FeatureWithMA150GTMA200})
	domain.RegisterFeatureFunc(&domain.Feature{Desc: "150 日均线向上延伸", F: FeatureWithMA150UP})
	domain.RegisterFeatureFunc(&domain.Feature{Desc: "150 日均线反转", F: FeatureWithReversalMA150})
	domain.RegisterFeatureFunc(&domain.Feature{Desc: "最近 5 天内收盘价靠近 150 日均线", F: FeatureWithNearbyMA150})
}

func FeatureWithClosedGTMA_50(k *chart.K) bool {
	if len(k.Candlesticks) == 0 {
		return false
	}

	lastCandletick := k.Candlesticks[len(k.Candlesticks)-1]
	ma_50, ok := lastCandletick.MA[chart.Ma_50]
	if !ok {
		return false
	}
	return lastCandletick.Close > ma_50
}

func FeatureWithClosedGTMA150(k *chart.K) bool {
	if len(k.Candlesticks) == 0 {
		return false
	}

	lastCandletick := k.Candlesticks[len(k.Candlesticks)-1]
	ma150, ok := lastCandletick.MA[chart.Ma150]
	if !ok {
		return false
	}
	return lastCandletick.Close > ma150
}

func FeatureWithClosedGTMA200(k *chart.K) bool {
	if len(k.Candlesticks) == 0 {
		return false
	}

	lastCandletick := k.Candlesticks[len(k.Candlesticks)-1]
	ma200, ok := lastCandletick.MA[chart.Ma200]
	if !ok {
		return false
	}
	return lastCandletick.Close > ma200
}

func FeatureWithMA150GTMA200(k *chart.K) bool {
	if len(k.Candlesticks) == 0 {
		return false
	}

	lastCandletick := k.Candlesticks[len(k.Candlesticks)-1]
	ma150, ok := lastCandletick.MA[chart.Ma150]
	if !ok {
		return false
	}
	ma200, ok := lastCandletick.MA[chart.Ma200]
	if !ok {
		return false
	}

	return ma150 > ma200
}

func FeatureWithMA150UP(k *chart.K) bool {
	if len(k.Candlesticks) <= 1 {
		return false
	}

	font, back := k.Candlesticks[len(k.Candlesticks)-2], k.Candlesticks[len(k.Candlesticks)-1]
	fontMA150, ok := font.MA[chart.Ma150]
	if !ok {
		return false
	}
	backMA150, ok := back.MA[chart.Ma150]
	if !ok {
		return false
	}
	return backMA150 > fontMA150
}

func FeatureWithReversalMA150(k *chart.K) bool {
	trend := calculateTrendWithMA(chart.Ma150, k)

	if len(trend) == 2 {
		s1, s2 := trend[len(trend)-2], trend[len(trend)-1]
		if len(s1) > 1 && (s1[len(s1)-1] <= s1[0]) && len(s2) > 1 && (s2[len(s2)-1] > s2[0]) {
			return true
		}
	}
	return false
}

func FeatureWithNearbyMA150(k *chart.K) bool {
	if len(k.Candlesticks) < 5 {
		return false
	}

	// for i := len(k.Candlesticks) - 5; i < len(k.Candlesticks); i++ {
	current := k.Candlesticks[len(k.Candlesticks)-1]

	closed := current.Close
	ma, ok := current.MA[chart.Ma150]
	if !ok {
		return false
	}

	if closed*0.97 < ma {
		return true
	}
	// }
	return false
}

func calculateTrendWithMA(kind chart.MaKind, k *chart.K) [][]float64 {
	if len(k.Candlesticks) == 0 {
		return nil
	}

	ma := make([]float64, 0, len(k.Candlesticks))
	begin := -1
	for i, candlestick := range k.Candlesticks {
		val, ok := candlestick.MA[kind]
		if !ok {
			continue
		}
		if begin == -1 {
			begin = i
		}
		ma = append(ma, val)
	}

	if len(ma) == 0 {
		return nil
	}

	trend := make([][]float64, 0, 4)

	for i := 0; i < len(ma); i++ {
		span := make([]float64, 0, 32)
		span = append(span, ma[i])

		direction := 0
	loop:
		for j := i + 1; j < len(ma); j++ {
			switch {
			case ma[i] > ma[j]:
				if direction == 0 {
					direction = 1
				}
				if direction == 1 {
					i = j
				}
				if direction == 2 {
					break loop
				}

			case ma[i] < ma[j]:
				if direction == 0 {
					direction = 2
				}
				if direction == 2 {
					i = j
				}
				if direction == 1 {
					break loop
				}
			default:
				i = j
			}
			span = append(span, ma[j])
		}
		trend = append(trend, span)

	}

	return trend
}
