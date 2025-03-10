package chart

import "github.com/eviltomorrow/king/lib/mathutil"

func CalculateMaToSegment(k *K, day int) [][]float64 {
	if len(k.Candlesticks) == 0 {
		return nil
	}

	ma := make([]float64, 0, len(k.Candlesticks))
	begin := -1
	for i, candlestick := range k.Candlesticks {
		val, ok := candlestick.Indicators.Trend.MA[day]
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

func CalculateMaOnNext(k *K, day, m int) ([]float64, error) {
	result := make([]float64, 0, day)

	x := make([]float64, 0, len(k.Candlesticks))
	y := make([]float64, 0, len(k.Candlesticks))
	n := 0
	for _, c := range k.Candlesticks {
		val, ok := c.Indicators.Trend.MA[day]
		if ok {
			n++
			x = append(x, float64(n))
			y = append(y, val)
		}

	}

	a, b, err := mathutil.LeastSquares(x, y)
	if err != nil {
		return nil, err
	}
	next := a*float64(n+1) + b
	result = append(result, mathutil.Trunc4(next))

	for i := 1; i < day; i++ {
		n = n + i
		x = append(x, float64(n))
		y = append(y, next)

		a, b, err := mathutil.LeastSquares(x, y)
		if err != nil {
			return nil, err
		}

		next = a*float64(n+1) + b
		result = append(result, next)

	}

	return result, nil
}
