package data

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	DAY  = "day"
	WEEK = "week"
)

func FetchStock(ctx context.Context, pipe chan *Stock) error {
	if pipe == nil {
		return fmt.Errorf("panic: invalid pipe")
	}
	defer func() {
		close(pipe)
	}()

	resp, err := client.DefalutStorage.GetStockAll(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	var e error = nil
	for {
		stock, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			e = err
			break
		}
		pipe <- &Stock{
			Code:    stock.Code,
			Name:    stock.Name,
			Suspend: stock.Suspend,
		}
	}

	return e
}

func GetQuote(ctx context.Context, date time.Time, code string, kind string) ([]*Quote, error) {
	reverse := func(quotes []*Quote) []*Quote {
		for i, j := 0, len(quotes)-1; i < j; i, j = i+1, j-1 {
			quotes[i], quotes[j] = quotes[j], quotes[i]
		}
		return quotes
	}

	var limit int64 = 250
	resp, err := client.DefalutStorage.GetQuoteLatest(ctx, &pb.GetQuoteLatestRequest{
		Code:  code,
		Date:  date.Format(time.DateOnly),
		Limit: limit,
		Kind: func() pb.GetQuoteLatestRequest_Kind {
			switch kind {
			case "day":
				return pb.GetQuoteLatestRequest_Day
			case "week":
				return pb.GetQuoteLatestRequest_Week
			default:
				return pb.GetQuoteLatestRequest_Day
			}
		}(),
	})
	if err != nil {
		return nil, err
	}

	data := make([]*Quote, 0, limit)
	for {
		quote, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		q := &Quote{
			Code:            quote.Code,
			Open:            quote.Open,
			Close:           quote.Close,
			High:            quote.High,
			Low:             quote.Low,
			YesterdayClosed: quote.YesterdayClosed,
			Volume:          quote.Volume,
			Account:         quote.Account,
			Date:            quote.Date,
			NumOfYear:       quote.NumOfYear,
		}
		data = append(data, q)
	}

	return reverse(data), nil
}
