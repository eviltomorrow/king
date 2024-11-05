package model

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/service"
	"github.com/eviltomorrow/king/lib/mathutil"
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

	score := 0
	var (
		lastCandlestick = k.Candlesticks[len(k.Candlesticks)-1]
		ma_10           = lastCandlestick.MA[chart.Ma_10]
		ma_50           = lastCandlestick.MA[chart.Ma_50]
		ma150           = lastCandlestick.MA[chart.Ma150]
		ma200           = lastCandlestick.MA[chart.Ma200]

		closed    = lastCandlestick.Close
		closed_95 = mathutil.Trunc2(closed * 0.95)
	)

	if ma_10 > ma_50 {
		score += SCORE_5
	}
	if ma150 > ma200 {
		score += SCORE_20
	}
	if closed > ma150 && ma150 > closed_95 {
		score += SCORE_50
	}

	return score, nil
}
