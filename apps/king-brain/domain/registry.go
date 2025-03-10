package domain

import "github.com/eviltomorrow/king/apps/king-brain/domain/chart"

var cache []*Model

type Model struct {
	Desc string
	F    func(k *chart.K) (*Plan, error)
}

func RegisterModel(model *Model) {
	cache = append(cache, model)
}

func ScanModel(k *chart.K) ([]*Plan, error) {
	plans := make([]*Plan, 0, len(cache))
	for _, model := range cache {
		plan, err := model.F(k)
		if err != nil {
			return nil, err
		}
		if plan == nil {
			continue
		}
		plans = append(plans, plan)
	}
	return plans, nil
}
