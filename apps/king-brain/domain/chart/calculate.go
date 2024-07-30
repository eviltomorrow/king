package chart

import (
	"github.com/eviltomorrow/king/lib/mathutil"
)

func calculateMa(closed []float64) float64 {
	sum := mathutil.Sum(closed)
	return mathutil.TruncN(sum/float64(len(closed)), 3)
}
