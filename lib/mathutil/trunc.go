package mathutil

import (
	"math"
)

var (
	n10_1 = math.Pow10(2)
	n10_4 = math.Pow10(4)
)

// Trunc2 trunc2
func Trunc2(val float64) float64 {
	return math.Trunc((val+0.5/n10_1)*n10_1) / n10_1
}

func Trunc4(val float64) float64 {
	return math.Trunc((val+0.5/n10_4)*n10_4) / n10_4
}

func TruncN(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
