package mathutil

import (
	"fmt"
	"testing"
)

func TestSumFloat64(t *testing.T) {
	sum := Sum([]float64{1, 2, 3, 4, 5})
	fmt.Println(sum)
}
