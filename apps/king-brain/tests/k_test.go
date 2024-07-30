package tests

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
)

func TestNewK(t *testing.T) {
	var pipe = make(chan *pb.Stock, 64)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for stock := range pipe {
			quotes, err := data.FetchQuote(context.Background(), time.Now(), stock.Code)
			if err != nil {
				log.Fatal(err)
			}

			k, err := chart.NewK(context.Background(), stock, quotes)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(k)
		}
		wg.Done()
	}()
	if err := data.FetchStock(context.Background(), pipe); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
