package data

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
	jsoniter "github.com/json-iterator/go"
)

type Quote struct {
	Code            string  `json:"code"`
	Open            float64 `json:"opend"`
	Close           float64 `json:"close"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	YesterdayClosed float64 `json:"yesterday_closed"`
	Volume          int64   `json:"volume"`
	Account         float64 `json:"account"`
	Date            string  `json:"date"`
	NumOfYear       int32   `json:"num_of_year"`
}

func (q *Quote) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	if err != nil {
		return fmt.Sprintf("marshal quote failure, nest error: %v", err)
	}
	return string(buf)
}

func GetQuotesN(ctx context.Context, date time.Time, code string, kind string, n int64) ([]*Quote, error) {
	reverse := func(quotes []*Quote) []*Quote {
		for i, j := 0, len(quotes)-1; i < j; i, j = i+1, j-1 {
			quotes[i], quotes[j] = quotes[j], quotes[i]
		}
		return quotes
	}

	var limit int64 = n
	resp, err := client.DefaultStorage.GetQuoteLatest(ctx, &pb.GetQuoteLatestRequest{
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
