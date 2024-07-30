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

func FetchStock(ctx context.Context, pipe chan *pb.Stock) error {
	if pipe == nil {
		return fmt.Errorf("panic: invalid pipe")
	}
	stub, closeFunc, err := client.NewStorageWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	resp, err := stub.GetStockAll(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	for {
		stock, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		pipe <- stock
	}
	close(pipe)

	return nil
}

func FetchQuote(ctx context.Context, date time.Time, code string) ([]*pb.Quote, error) {
	reverse := func(quotes []*pb.Quote) []*pb.Quote {
		for i, j := 0, len(quotes)-1; i < j; i, j = i+1, j-1 {
			quotes[i], quotes[j] = quotes[j], quotes[i]
		}
		return quotes
	}

	stub, closeFunc, err := client.NewStorageWithEtcd()
	if err != nil {
		return nil, err
	}
	defer closeFunc()

	var limit int64 = 250
	resp, err := stub.GetQuoteLatest(ctx, &pb.GetQuoteLatestRequest{
		Code:  code,
		Date:  date.Format(time.DateOnly),
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	var data = make([]*pb.Quote, 0, limit)
	for {
		quote, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		data = append(data, quote)
	}

	return reverse(data), nil
}
