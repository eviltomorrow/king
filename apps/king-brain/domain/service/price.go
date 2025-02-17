package service

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
)

type Plan struct{}

func AnalysisChartFindBuyPoint(ctx context.Context, k *chart.K) (*Plan, bool) {
	if len(k.Candlesticks) == 0 {
		return nil, false
	}

	for _, c := range k.Candlesticks {
		fmt.Println(c.String())
	}
	return nil, false
}

func assume()
