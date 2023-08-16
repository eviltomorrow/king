package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataWrapperChannel(t *testing.T) {
	assert := assert.New(t)
	pipe, err := NewDataWrapperChannel(DAY, "2023-08-03")
	assert.Nil(err)

	var count int
	for d := range pipe {
		if len(d.Data) > 1 {
			var (
				pre  = d.Data[len(d.Data)-2]
				next = d.Data[len(d.Data)-1]
			)
			if next.Quote.Close > next.MA_50 && next.MA_50 > next.MA150 && next.MA150 > next.MA200 && next.MA200 > pre.MA200 && next.MA_50 > pre.MA_50 {
				count++
				t.Logf("%s\r\n", next.String())
			}
		}
	}
	fmt.Println(count)
}
