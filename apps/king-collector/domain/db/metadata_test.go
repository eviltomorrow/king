package db

import (
	"context"
	"testing"

	"github.com/eviltomorrow/king/lib/db/mongodb"
	"github.com/stretchr/testify/assert"
)

func TestSelectMetadataRange(t *testing.T) {
	assert := assert.New(t)

	data, err := SelectMetadataRange(context.Background(), mongodb.DB, 0, 30, "2024-10-11", "")
	assert.Nil(err)

	for _, d := range data {
		t.Logf("data: %v", d.String())
	}
}

func TestSelectMetadataAll(t *testing.T) {
	var (
		offset, limit int64 = 0, 30
		lastID        string
		count         int
	)

	for {
		metadata, err := SelectMetadataRange(context.Background(), mongodb.DB, offset, limit, "2024-10-11", lastID)
		if err != nil {
			t.Fatal(err)
		}

		for _, md := range metadata {
			count++
			_ = md
		}
		if len(metadata) < int(limit) {
			break
		}
		offset += limit
	}
	t.Logf("count: %v", count)
}
