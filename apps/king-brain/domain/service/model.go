package service

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	jsoniter "github.com/json-iterator/go"
)

type Position struct {
	Buy      float64
	StopLoss float64
}

func (p *Position) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(p)
	return string(buf)
}

type Model struct {
	Name string
	Desc string

	F func(*chart.K) (*Position, bool)
}

var repository = make([]*Model, 0, 8)

func RegisterModel(m *Model) {
	repository = append(repository, m)
}

func ScanModel(k *chart.K) (int, bool, error) {
	sum := 0
	ok := true
	for _, m := range repository {
		position, ok := m.F(k)
		if ok {
			fmt.Println(position)
		}

	}
	return sum, ok, nil
}
