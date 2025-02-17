package service

import (
	"context"
	"testing"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	"github.com/stretchr/testify/assert"
)

func TestAnalysisChart(t *testing.T) {
	assert := assert.New(t)

	date := time.Now()
	stock := &data.Stock{
		Code: "sh601933",
		Name: "永辉超市",
	}

	ctx := context.Background()

	quotes, err := data.GetQuotesN(ctx, date, stock.Code, "day", 250)
	assert.Nil(err)

	k, err := chart.NewK(ctx, stock, quotes)
	assert.Nil(err)

	AnalysisChartFindBuyPoint(ctx, k)
}
