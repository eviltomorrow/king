package model

import (
	"fmt"
	"math/big"

	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/lib/mathutil"
)

func init() {
	domain.RegisterModel(&domain.Model{Desc: "", F: F_01})
}

func F_01(k *chart.K) (*domain.Plan, bool) {
	if len(k.Candlesticks) <= 200 {
		return nil, false
	}

	k.CalMoreMA([]int{150, 200})

	fmt.Println(k.String())

	return &domain.Plan{K: k}, true
}

func InferClosedbyMa(k *chart.K, day int) (float64, bool) {
	if len(k.Candlesticks) <= day || day < 3 {
		return 0, false
	}

	k.CalMA(day)

	var (
		ma     = make([]float64, 0, len(k.Candlesticks))
		closed = make([]float64, day)
	)

	n := 0
	for _, c := range k.Candlesticks {
		if n <= day-1 {
			closed[n] = c.Close
			n++
		} else {
			for i := 0; i < n-2; i++ {
				closed[i], closed[i+1] = closed[i+1], closed[i+2]
			}
			closed[len(closed)-1] = c.Close
		}

		val, ok := c.Indicators.Trend.MA[day]
		if ok {
			ma = append(ma, val)
		}
	}

	fmt.Println(k.Candlesticks[len(k.Candlesticks)-1])
	// c, err := InferNextMa(ma, day)
	// if err != nil {
	// 	zlog.Error("infer ma failure", zap.Error(err))
	// 	return 0, false
	// }
	// fmt.Println(c)
	// // c = 750.44

	c := 4.28
	// fmt.Println(closed)
	// sum1 := mathutil.Sum(closed[:])
	sum2 := mathutil.SumFloat64(closed[1:])
	dd, _ := sum2.Float64()
	aa := new(big.Float).SetPrec(64).SetFloat64(c)
	bb := aa.Mul(aa, new(big.Float).SetPrec(64).SetFloat64(150))
	cc := bb.Sub(bb, new(big.Float).SetPrec(64).SetFloat64(dd))
	fmt.Println(cc)
	// fmt.Println(sum1, sum2)
	// fmt.Println((sum2 + 5.02) / 150)
	// fmt.Println(sum1/150, sum2)
	// fmt.Println(c*float64(day) - sum)

	return 0, true
}

func InferNextMa(ma []float64, day int) (float64, error) {
	y := ma

	if len(y) > day {
		y = y[len(y)-day:]
	}

	x := make([]float64, 0, len(y))
	for i := 1; i <= len(y); i++ {
		x = append(x, float64(i))
	}

	fmt.Println(len(x), x)
	fmt.Println(y)
	a, b, err := mathutil.LeastSquares(x, y)
	if err != nil {
		return 0, err
	}

	return a*float64(len(x)+1) + b, nil
}
