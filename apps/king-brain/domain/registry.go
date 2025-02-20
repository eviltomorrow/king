package domain

import "github.com/eviltomorrow/king/apps/king-brain/domain/chart"

var cache []*Model

type Model struct {
	Desc string
	F    func(k *chart.K) (*Plan, bool)
}

func RegisterModel(model *Model) {
	cache = append(cache, model)
}

func ScanModel(k *chart.K) []*Plan {
	plans := make([]*Plan, 0, len(cache))
	for _, model := range cache {
		plan, ok := model.F(k)
		if ok {
			plans = append(plans, plan)
		}
	}
	return plans
}
