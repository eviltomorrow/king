package model

import (
	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func init() {
	domain.RegisterFeaturesFunc(&domain.Feature{Name: "收盘价高于 50 日均线", F: FeatureWithClosedGTMA_50})
	domain.RegisterFeaturesFunc(&domain.Feature{Name: "收盘价高于 150 日均线", F: FeatureWithClosedGTMA150})
	domain.RegisterFeaturesFunc(&domain.Feature{Name: "收盘价高于 200 日均线", F: FeatureWithClosedGTMA150})
	domain.RegisterFeaturesFunc(&domain.Feature{Name: "150 日均线高于 200 日均线", F: FeatureWithMA150GTMA200})
}

func FeatureWithClosedGTMA_50(k *chart.K) (int64, bool) {
	if len(k.Candlesticks) == 0 {
		return 0, false
	}

	lastCandletick := k.Candlesticks[len(k.Candlesticks)-1]
	ma_50, ok := lastCandletick.MA[chart.Ma_50]
	if !ok {
		return 0, false
	}
	if lastCandletick.Close > ma_50 {
		return 1, true
	}
	return 0, false
}

func FeatureWithClosedGTMA150(k *chart.K) (int64, bool) {
	if len(k.Candlesticks) == 0 {
		return 0, false
	}

	lastCandletick := k.Candlesticks[len(k.Candlesticks)-1]
	ma150, ok := lastCandletick.MA[chart.Ma150]
	if !ok {
		return 0, false
	}
	if lastCandletick.Close > ma150 {
		return 2, true
	}
	return 0, false
}

func FeatureWithClosedGTMA200(k *chart.K) (int64, bool) {
	if len(k.Candlesticks) == 0 {
		return 0, false
	}

	lastCandletick := k.Candlesticks[len(k.Candlesticks)-1]
	ma200, ok := lastCandletick.MA[chart.Ma200]
	if !ok {
		return 0, false
	}
	if lastCandletick.Close > ma200 {
		return 2, true
	}
	return 0, false
}

func FeatureWithMA150GTMA200(k *chart.K) (int64, bool) {
	if len(k.Candlesticks) == 0 {
		return 0, false
	}

	lastCandletick := k.Candlesticks[len(k.Candlesticks)-1]
	ma150, ok := lastCandletick.MA[chart.Ma150]
	if !ok {
		return 0, false
	}
	ma200, ok := lastCandletick.MA[chart.Ma200]
	if !ok {
		return 0, false
	}
	if ma150 > ma200 {
		return 3, true
	}

	return 0, false
}
