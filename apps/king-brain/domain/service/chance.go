package service

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func FindPossibleChance(k *chart.K) {
	ScanFeatureFunc(k)
}

var features = make([]*Feature, 0, 32)

type Feature struct {
	Name string
	F    func(*chart.K) (int64, bool)
}

func RegisterFeaturesFunc(f *Feature) {
	if f != nil {
		features = append(features, f)
	}
}

func ScanFeatureFunc(k *chart.K) {
	for _, feature := range features {
		score, ok := feature.F(k)
		if !ok {
			continue
		}
		fmt.Printf("name: %s, score: %v\r\n", feature.Name, score)
	}
}
