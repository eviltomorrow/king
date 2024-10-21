package mathutil

import "math"

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func Angle(m Point, n Point) (angel float64) {
	if n.X <= m.X {
		return -180.0
	}
	return Trunc2(math.Atan2(n.Y-m.Y, n.X-m.X) * (180 / math.Pi))
}
