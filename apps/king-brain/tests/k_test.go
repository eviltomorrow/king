package tests

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
)

func TestNewK(t *testing.T) {
	quotes, err := data.GetQuotesN(context.Background(), time.Now(), "sh600519", "day", 2)
	if err != nil {
		log.Fatal(err)
	}

	k, err := chart.NewK(context.Background(), &data.Stock{Name: "贵州茅台", Code: "sh600519"}, quotes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k)
}

func TestFind(t *testing.T) {
	pipe := make(chan *data.Stock, 64)

	go func() {
		if err := data.FetchStock(context.Background(), pipe); err != nil {
			log.Fatal(err)
		}
	}()

	var wg sync.WaitGroup
	var count atomic.Int32
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			for stock := range pipe {
				quotes, err := data.GetQuotesN(context.Background(), time.Now(), stock.Code, "day", 3)
				if err != nil {
					log.Printf("GetQuote failure, nest error: %v", err)
				} else {
					k, err := chart.NewK(context.Background(), stock, quotes)
					if err != nil {
						log.Printf("New k failure, nest error: %v", err)
					}
					_ = k
					count.Add(1)
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(count.Load())
}
