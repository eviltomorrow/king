package model

import (
	"fmt"
	"log"

	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/lib/mathutil"
)

func init() {
	domain.RegisterModel(&domain.Model{Desc: "", F: F_01})
}

func F_01(k *chart.K) (*domain.Plan, error) {
	if len(k.Candlesticks) <= 200 {
		return nil, nil
	}

	days := []int{10}
	k.CalMaMany(days)

	data, err := chart.CalculateMaToSegment(k, 10)
	if err != nil {
		log.Fatal(err)
	}

	next, err := chart.CalculateMaOnNext(k, 10, 10, 3)
	if err != nil {
		log.Fatal(err)
	}

	next, err = forecastPriveOnNext(data[len(data)-1], 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(next)

	fmt.Println(next)
	return &domain.Plan{K: k}, nil
}

// func forecastCandlesticskOnNext(k *chart.K, span, count int) ([]*chart.Candlestick, error) {
// 	if len(k.Candlesticks) <= span {
// 		return nil, fmt.Errorf("no enough data")
// 	}

// 	var (
// 		high = make([]float64, 0, span+count)
// 		low  = make([]float64, 0, span+count)

// 		open   = make([]float64, 0, span+count)
// 		closed = make([]float64, 0, span+count)
// 	)

// 	for i := len(k.Candlesticks) - span; i < len(k.Candlesticks); i++ {
// 		c := k.Candlesticks[i]

// 		high = append(high, c.High*100)
// 		low = append(low, c.Low*100)
// 		open = append(open, c.Open*100)
// 		closed = append(closed, c.Close)
// 	}

// 	data, err := chart.CalculateMaToSegment(k, 150)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(data)
// 	fClosed, err := forecastPriveOnNext(data[len(data)-1], count)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(fClosed)
// 	return nil, nil
// }

func forecastPriveOnNext(data []float64, count int) ([]float64, error) {
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
