package feature

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func init() {
	domain.RegisterFeatureFunc(&domain.Feature{Level: Level2, Desc: "收盘价高于 150 日均线", F: FeatureWithClosedGTMA150})
	domain.RegisterFeatureFunc(&domain.Feature{Level: Level2, Desc: "收盘价高于 200 日均线", F: FeatureWithClosedGTMA150})
	domain.RegisterFeatureFunc(&domain.Feature{Level: Level2, Desc: "150 日均线高于 200 日均线", F: FeatureWithMA150GTMA200})
	domain.RegisterFeatureFunc(&domain.Feature{Level: Level2, Desc: "150 日均线向上延伸", F: FeatureWithMA150UP})
	domain.RegisterFeatureFunc(&domain.Feature{Level: Level2, Desc: "150 日均线反转", F: FeatureWithMA150Reversal})
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

func FeatureWithMA150Reversal(k *chart.K) bool {
	if len(k.Candlesticks) == 0 {
		return false
	}

	ma150 := make([]float64, 0, len(k.Candlesticks))
	begin := -1
	for i, candlestick := range k.Candlesticks {
		val, ok := candlestick.MA[chart.Ma150]
		if !ok {
			continue
		}
		if begin == -1 {
			begin = i
		}
		ma150 = append(ma150, val)
	}

	if len(ma150) == 0 {
		return false
	}

	swing := make([][]float64, 0, 4)

	for i := 0; i < len(ma150); i++ {
		span := make([]float64, 0, 32)
		span = append(span, ma150[i])

		direction := 0
	loop:
		for j := i + 1; j < len(ma150); j++ {
			switch {
			case ma150[i] > ma150[j]:
				if direction == 0 {
					direction = 1
				}
				if direction == 1 {
					i = j
				}
				if direction == 2 {
					break loop
				}

			case ma150[i] < ma150[j]:
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
			span = append(span, ma150[j])
		}
		swing = append(swing, span)

	}

	if len(swing) == 2 {
		s1, s2 := swing[len(swing)-2], swing[len(swing)-1]
		if len(s1) > 1 && (s1[len(s1)-1] <= s1[0]) && len(s2) > 1 && (s2[len(s2)-1] > s2[0]) {
			return true
		}
	}
	return false
}
