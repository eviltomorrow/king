package model

import (
	"fmt"
	"log"

	"github.com/eviltomorrow/king/apps/king-brain/domain"
	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

func init() {
	domain.RegisterModel(&domain.Model{Desc: "", F: F_01})
}

func F_01(k *chart.K) (*domain.Plan, bool) {
	if len(k.Candlesticks) <= 200 {
		return nil, false
	}

	k.CalMa(150)

	data, err := chart.CalculateMaOnNext(k, 150, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)

	return &domain.Plan{K: k}, true
}
