package model

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/service"
)

type ConfigM1 struct{}

var C1 = ConfigM1{}

func init() {
	service.RegisterModel(&service.Model{
		Name: "",
		Desc: "",
		F:    M1,
	})
}

func M1(k *chart.K) (*service.Position, bool) {
	if len(k.Candlesticks) == 0 {
		return nil, false
	}

	return nil, false
}
