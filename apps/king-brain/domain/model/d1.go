package model

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/service"
)

func init() {
	service.RegisterModel(&service.Model{
		Name: "",
		Desc: "",
		F:    d1,
		C: func(score int) bool {
			if score >= 75 {
				return true
			}
			return false
		},
	})
}

func d1(k *chart.K) (int, error) {
	if len(k.Candlesticks) == 0 {
		return 0, nil
	}

	// score := 0
	// var (
	// 	lastCandlestick = k.Candlesticks[len(k.Candlesticks)-1]
	// 	ma_10           = lastCandlestick.MA[chart.Ma_10]
	// 	ma_50           = lastCandlestick.MA[chart.Ma_50]
	// 	ma150           = lastCandlestick.MA[chart.Ma150]
	// 	ma200           = lastCandlestick.MA[chart.Ma200]

	// 	closed    = lastCandlestick.Close
	// 	closed_95 = mathutil.Trunc2(closed * 0.95)
	// )

	return 0, nil
}
