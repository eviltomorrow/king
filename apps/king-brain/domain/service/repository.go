package service

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

var repository = map[string]*Model{}

type Model struct {
	Name string
	Desc string
	F    func(k *chart.K) (int, bool)
}

func RegisterModel(name string, model *Model) {
	repository[name] = model
}

func FindPossibleChance(k *chart.K) error {
	for _, model := range repository {
		score, ok := model.F(k)
		if !ok {
			continue
		}

		fmt.Println(k.Code, k.Name, score)

	}
	return nil
}
