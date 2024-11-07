package datasource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchMetadataFromSina(t *testing.T) {
	_assert := assert.New(t)
	codes := []string{
		"bj899050",
		"sh000688",
	}
	data, err := FetchMetadataFromSina(codes)
	_assert.Nil(err)
	t.Logf("data: %s\r\n", data)

	// [{"_id":"","source":"","code":"sh689009","name":"九号公司","open":44.61,"yesterday_closed":43.03,"latest":48.2,"high":49.24,"low":44.06,"volume":26790226,"account":1248638943,"date":"2024-09-30","time":"15:00:03","suspend":"正常"}]
}
