package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataWrapperChannel(t *testing.T) {
	assert := assert.New(t)
	pipe, err := NewDataWrapperChannel(DAY, "2023-06-27")
	assert.Nil(err)

	var count int
	for d := range pipe {
		if len(d.Data) > 0 {
			var data = d.Data[len(d.Data)-1]
			if data.Quote.Close > data.MA_50 && data.MA_50 > data.MA150 && data.MA150 > data.MA200 {
				count++
				t.Logf("%s %s\r\n", d.Stock.Code, d.Stock.Name)
			}
		}
	}
	fmt.Println(count)
}
