package service

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

type Model struct {
	Name string
	Desc string

	C func(int) bool
	F func(*chart.K) (int, error)
}

var repository = make([]*Model, 0, 8)

func RegisterModel(m *Model) {
	repository = append(repository, m)
}

func ScanModel(k *chart.K) (int, bool, error) {
	sum := 0
	ok := true
	for _, m := range repository {
		score, err := m.F(k)
		if err != nil {
			return 0, false, err
		}

		ok = m.C(score)
		if !ok {
			return 0, false, nil
		}
		sum += score
	}
	return sum, ok, nil
}
