package mathutil

import (
	"fmt"
	"testing"
)

func TestAngle(t *testing.T) {
	m := Point{
		X: 0,
		Y: 0,
	}
	n := Point{
		X: 1,
		Y: -0.3,
	}

	var data = Angle(m, n)
	fmt.Println(data)

}
