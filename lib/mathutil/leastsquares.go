package mathutil

import "fmt"

func LeastSquares(x []float64, y []float64) (a float64, b float64, err error) {
	xi := float64(0)
	x2 := float64(0)
	yi := float64(0)
	xy := float64(0)

	if len(x) != len(y) {
		return 0, 0, fmt.Errorf("[]x != []y")
	} else {
		length := float64(len(x))
		for i := 0; i < len(x); i++ {
			xi += x[i]
			x2 += x[i] * x[i]
			yi += y[i]
			xy += x[i] * y[i]
		}
		a = (yi*xi - xy*length) / (xi*xi - x2*length)
		b = (yi*x2 - xy*xi) / (x2*length - xi*xi)
	}
	return
}
