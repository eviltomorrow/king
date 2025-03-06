package mathutil

import (
	"math/big"
	"math/rand"
	"time"
)

func GenRandInt(min, max int) int {
	return rand.New(rand.NewSource(time.Now().Unix())).Intn(max-min) + min
}

func Max[T int | uint | int64 | uint64 | float64](data []T) T {
	if len(data) == 0 {
		return 0
	}
	max := data[0]
	for i := 1; i <= len(data)-1; i++ {
		if data[i] > max {
			max = data[i]
		}
	}
	return max
}

func Min[T int | uint | int64 | uint64 | float64](data []T) T {
	if len(data) == 0 {
		return 0
	}
	min := data[0]
	for i := 1; i <= len(data)-1; i++ {
		if data[i] < min {
			min = data[i]
		}
	}
	return min
}

func Sum[T int | uint | int64 | uint64 | float64](data []T) T {
	var sum T
	for _, d := range data {
		sum += d
	}
	return sum
}

func SumFloat64(data []float64) *big.Float {
	sum := &big.Float{}
	for _, d := range data {
		f := new(big.Float).SetPrec(64).SetFloat64(d)
		sum.Add(sum, f)
	}
	return sum
}
