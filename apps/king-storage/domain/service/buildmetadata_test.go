package service

import (
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/eviltomorrow/king/lib/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildQuoteDayWitchMetadata(t *testing.T) {
	gomonkey.ApplyFunc(BuildQuoteDayWitchMetadata)

	assert := assert.New(t)
	data := &model.Metadata{
		Source:          "sina",
		Code:            "sh601066",
		Name:            "中信建投",
		Open:            22.10,
		YesterdayClosed: 22.0,
		Latest:          23.0,
		High:            24.01,
		Low:             22.10,
		Volume:          108202301,
		Account:         2930241023,
		Date:            "2024-09-27",
		Time:            "15:00:03",
		Suspend:         "正常",
	}

	date, err := time.Parse(time.DateOnly, "2024-09-27")
	assert.Nil(err)

	quote, err := BuildQuoteDayWitchMetadata(data, date)
	assert.Nil(err)
	assert.NotNil(quote)
}
