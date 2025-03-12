package chart

import (
	"errors"

	"github.com/eviltomorrow/king/lib/mathutil"
)

var ErrNoData = errors.New("no data")

func CalculateMaToSegment(k *K, day int) ([][]float64, error) {
	if len(k.Candlesticks) == 0 {
		return nil, ErrNoData
	}

	data := make([]float64, 0, len(k.Candlesticks))
	begin := -1
	for i, candlestick := range k.Candlesticks {
		ma, ok := candlestick.Indicators.Trend.Ma[day]
		if !ok {
			continue
		}
		if begin == -1 {
			begin = i
		}
		data = append(data, ma)
	}

	if len(data) == 0 {
		return nil, ErrNoData
	}

	trend := make([][]float64, 0, 4)

	for i := 0; i < len(data); i++ {
		span := make([]float64, 0, 32)
		span = append(span, data[i])

		direction := 0
	loop:
		for j := i + 1; j < len(data); j++ {
			switch {
			case data[i] > data[j]:
				if direction == 0 {
					direction = 1
				}
				if direction == 1 {
					i = j
				}
				if direction == 2 {
					break loop
				}

			case data[i] < data[j]:
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
			span = append(span, data[j])
		}
		trend = append(trend, span)
	}

	return trend, nil
}

func CalculateMaOnNext(k *K, day, span, count int) ([]float64, error) {
	result := make([]float64, 0, count)

	x := make([]float64, 0, len(k.Candlesticks)+count)
	y := make([]float64, 0, len(k.Candlesticks)+count)
	n := 0
	for _, c := range k.Candlesticks {
		ma, ok := c.Indicators.Trend.Ma[day]
		if ok {
			n++
			x = append(x, float64(n))
			y = append(y, ma)
		}
	}

	if len(x) == 0 || len(y) == 0 {
		return nil, ErrNoData
	}

	if len(x) > span {
		x = x[len(x)-span:]
	}
	if len(y) > span {
		y = y[len(y)-span:]
	}

	a, b, err := mathutil.LeastSquares(x, y)
	if err != nil {
		return nil, err
	}
	next := a*float64(n+1) + b
	result = append(result, mathutil.Trunc4(next))

	for i := 1; i < count; i++ {
		n = n + i
		x = append(x, float64(n))
		y = append(y, next)

		if len(x) > span {
			x = x[len(x)-span:]
		}
		if len(y) > span {
			y = y[len(y)-span:]
		}

		a, b, err = mathutil.LeastSquares(x, y)
		if err != nil {
			return nil, err
		}

		next = a*float64(n+1) + b
		result = append(result, mathutil.Trunc4(next))
	}
	return result, nil
}

func CalculateClosedOnNext(k *K, day, span int, count int) ([]float64, error) {
	// if len(k.Candlesticks) < span+1 {
	// 	return nil, fmt.Errorf("no enough data")
	// }

	// closed := make([]float64, 0, len(k.Candlesticks)+len(mas))
	// for _, c := range k.Candlesticks {
	// 	closed = append(closed, c.Close)
	// }

	result := make([]float64, 0, span)

	// sum := 0.0
	// for i := 0; i < len(mas); i++ {
	// 	sum = mathutil.Sum(closed[len(closed)-span+1:])
	// 	tmp := mas[i] * float64(span)

	// 	next := tmp - sum
	// 	closed = append(closed, next)
	// 	result = append(result, next)
	// }

	return result, nil
}
