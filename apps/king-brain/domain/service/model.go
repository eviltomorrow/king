package service

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

type Chance struct{}

type Model struct {
	Name string
	Desc string

	C func(int) bool
	F func(*chart.K) (*Chance, error)
}

var repository = make([]*Model, 0, 8)

func RegisterModel(m *Model) {
	repository = append(repository, m)
}

func ScanModel(k *chart.K) (int, bool, error) {
	sum := 0
	ok := true
	for _, m := range repository {
		chance, err := m.F(k)
		if err != nil {
			return 0, false, err
		}
		_ = chance

		ok = m.C(0)
		if !ok {
			return 0, false, nil
		}

	}
	return sum, ok, nil
}
