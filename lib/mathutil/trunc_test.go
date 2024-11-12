package mathutil

import (
	"fmt"
	"testing"
)

func TestTrunc2(t *testing.T) {
	c := Trunc2(float64(20.3-2.33)/float64(20.35)) * 100
	fmt.Println(c)
}
