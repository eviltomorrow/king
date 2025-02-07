package domain

import (
	"strings"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

var features = make([]*Feature, 0, 32)

type Feature struct {
	Desc string

	F func(k *chart.K) bool
}

func RegisterFeatureFunc(f *Feature) {
	if f != nil {
		features = append(features, f)
	}
}

func CalculateFeatureFunc(k *chart.K) (string, int64) {
	var (
		sum  int64
		desc = make([]string, 0, len(features))
	)
	for _, feature := range features {
		ok := feature.F(k)
		if !ok {
			continue
		}
		sum += 1
		desc = append(desc, feature.Desc)
	}
	return strings.Join(desc, ","), sum
}
