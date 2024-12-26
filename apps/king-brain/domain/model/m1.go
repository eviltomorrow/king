package model

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/service"
)

func init() {
	service.RegisterModel(&service.Model{
		Name: "",
		Desc: "",
		F:    m1,
		C: func(score int) bool {
			if score >= 75 {
				return true
			}
			return false
		},
	})
}

func m1(k *chart.K) (int, error) {
	if len(k.Candlesticks) == 0 {
		return 0, nil
	}

	return 0, nil
}
