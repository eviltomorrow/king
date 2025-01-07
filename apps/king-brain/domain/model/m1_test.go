package model

import (
	"context"
	"testing"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	"github.com/stretchr/testify/assert"
)

func TestM1(t *testing.T) {
	assert := assert.New(t)

	current, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	assert.Nil(err)

	stock := &data.Stock{
		Code: "sh601933",
		Name: "永辉超市",
	}
	quotes, err := data.GetQuotesN(context.Background(), current, stock.Code, "day", 250)
	assert.Nil(err)

	k, err := chart.NewK(context.Background(), stock, quotes)
	assert.Nil(err)

	position, ok := M1(k)
	if ok {
		t.Log(position)
	}
}
