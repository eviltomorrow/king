package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataWrapperChannel(t *testing.T) {
	assert := assert.New(t)
	pipe, err := NewDataWrapperChannel(DAY, "2023-11-20")
	assert.Nil(err)

	var count int
	for d := range pipe {
		if len(d.Data) > 1 {
			next := d.Data[len(d.Data)-1]
			if next.Quote.Close > next.MA_10 {
				count++
				t.Logf("%s\r\n", next.String())
			}
		}
	}
	fmt.Println(count)
}
