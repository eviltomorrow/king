package service

import (
	"context"
	"time"
)

type StatsInfo struct {
	Date string
	Kind string

	Desc map[string]string
}

func Report(ctx context.Context, date time.Time, kind string) (*StatsInfo, error) {
	// var (
	// 	wg sync.WaitGroup

	// 	pipe   = make(chan *data.Stock, 64)
	// 	result = make(chan *StatsInfo, 64)
	// )

	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		for stock := range pipe {
	// 			quotes, err := data.GetQuote(ctx, date, stock.Code, kind)
	// 			if err != nil {
	// 				zlog.Error("GetQuote failure", zap.Error(err), zap.String("code", stock.Code))
	// 				continue
	// 			}
	// 		}

	// 		wg.Done()
	// 	}()
	// }

	// go func() {
	// 	if err := data.FetchStock(context.Background(), pipe); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// wg.Wait()

	return nil, nil
}
