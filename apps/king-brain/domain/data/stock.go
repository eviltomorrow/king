package data

import (
	"context"
	"fmt"
	"io"

	"github.com/eviltomorrow/king/lib/grpc/client"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Stock struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Suspend string `json:"suspend"`
}

func (s *Stock) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(s)
	if err != nil {
		return fmt.Sprintf("marshal stock failure, nest error: %v", err)
	}
	return string(buf)
}

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
