package model

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/service"
)

type ConfigM1 struct {
	SlidingWindows int
}

var C1 = ConfigM1{
	SlidingWindows: 10,
}

func init() {
	service.RegisterModel(&service.Model{
		Name: "",
		Desc: "",
		F:    M1,
	})
}

func M1(k *chart.K) (*service.Strategy, bool) {
	if len(k.Candlesticks) == 0 {
		return nil, false
	}

	for i := len(k.Candlesticks) - 1; i >= 0; i-- {
	}
	return nil, false
}
