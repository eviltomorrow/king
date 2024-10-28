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
	quotes, err := data.GetQuote(context.Background(), time.Now(), "sz300004", "day")
	if err != nil {
		log.Fatal(err)
	}

	k, err := chart.NewK(context.Background(), &data.Stock{Name: "贵州茅台", Code: "sz300004"}, quotes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k)
}
