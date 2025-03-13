package chart

import (
	"errors"
	"fmt"

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

func CalculateMaOnNext(k *K, day, count int) ([]float64, error) {
	segments, err := CalculateMaToSegment(k, day)
	if err != nil {
		return nil, err
	}
	if len(segments) == 0 {
		return nil, fmt.Errorf("no ma segment")
	}

	segment := segments[len(segments)-1]

	if len(segment) > day {
		segment = segment[len(segment)-day:]
	}

	return CalculateOnNext(segment, count)
}

func CalculateOnNext(data []float64, count int) ([]float64, error) {
	result := make([]float64, 0, count)

	x := make([]float64, 0, len(data)+count)
	y := make([]float64, 0, len(data)+count)
	n := 0

	for _, d := range data {
		n++
		x = append(x, float64(n))
		y = append(y, d)
	}

	a, b, err := mathutil.LeastSquares(x, y)
	if err != nil {
		return nil, err
	}

	next := a*float64(n+1) + b
	result = append(result, mathutil.Trunc4(next))

	for i := 1; i < count; i++ {
		y = append(y, next)
		y = y[1:]

		a, b, err = mathutil.LeastSquares(x, y)
		if err != nil {
			return nil, err
		}

		next = a*float64(n+1) + b
		result = append(result, mathutil.Trunc4(next))
	}
	return result, nil
}
