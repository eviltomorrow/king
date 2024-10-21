package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
)

func TestNewK(t *testing.T) {
	quotes, err := data.FetchQuote(context.Background(), time.Now(), "sh600519", "day")
	if err != nil {
		log.Fatal(err)
	}

	k, err := chart.NewK(context.Background(), &data.Stock{Name: "贵州茅台", Code: "sh600519"}, quotes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k)
}
