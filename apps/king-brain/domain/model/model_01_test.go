package model

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	"github.com/stretchr/testify/assert"
)

func TestF01(t *testing.T) {
	assert := assert.New(t)

	stock := &data.Stock{
		Code: "sh601933",
		// Code: "sh688256",
		Name: "--",
	}
	date := time.Date(2025, time.March, 10, 12, 0, 0, 0, time.Local)

	quotes, err := data.GetQuotesN(context.Background(), date, stock.Code, "day", 250)
	assert.Nil(err)

	k, err := chart.NewK(context.Background(), stock, quotes)
	assert.Nil(err)

	plan, err := F_01(k)
	if err != nil {
		log.Fatal(err)
	}
	_ = plan
}
