package metadata

import (
	"context"
	"testing"
)

func TestSynchronizeMetadataQuick(t *testing.T) {
	total, ignore, err := SynchronizeMetadataQuick(context.Background(), "sina", []string{
		"bj920118",
		"bj920099",
		"bj920016",
		"bj920008",
		"bj920002",
		"bj873***",
		"bj872***",
		"bj871***",
		"bj83****",
		"sh689009",
		"sh688***",
		"sh605***",
		"sh603***",
		"sh601***",
		"sh600***",
		"sz301***",
		"sz300***",
		"sz003816",
		"sz0030**",
		"sz002***",
		"sz001***",
		"sz000***",

		"sh000001",
		"sz399001",
		"sz399006",
		"sh000688",
		"bj899050"},
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("total: %v, ignore: %v", total, ignore)
}
