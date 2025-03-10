package chart

import (
	"github.com/eviltomorrow/king/lib/mathutil"
)

func CalculateMa(closed []float64) float64 {
	sum := mathutil.Sum(closed)
	return sum / float64(len(closed))
}

func CalculateTrendWithMA(kind int, k *K) [][]float64 {
	if len(k.Candlesticks) == 0 {
		return nil
	}

	ma := make([]float64, 0, len(k.Candlesticks))
	begin := -1
	for i, candlestick := range k.Candlesticks {
		val, ok := candlestick.Indicators.Trend.MA[kind]
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
