package service

import (
	"context"
	"testing"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	"github.com/stretchr/testify/assert"
)

func TestFindPossibleChance(t *testing.T) {
	assert := assert.New(t)

	current, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	assert.Nil(err)

	stock := &data.Stock{
		Code: "sz301633",
		Name: "-",
	}
	quotes, err := data.GetQuotesN(context.Background(), current, stock.Code, "day", 250)
	assert.Nil(err)

	k, err := chart.NewK(context.Background(), stock, quotes)
	assert.Nil(err)

	FindPossibleChance(k)
}
