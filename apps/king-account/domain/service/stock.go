package service

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/lib/grpc/client"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Stock struct {
	Code string
	Name string
}

func StockWithGetOne(ctx context.Context, code string) (*Stock, error) {
	if code == "" {
		return nil, fmt.Errorf("code is nil")
	}
	stub, shutdown, err := client.NewStorageWithEtcd()
	if err != nil {
		return nil, err
	}
	defer shutdown()

	resp, err := stub.GetStockOne(ctx, &wrapperspb.StringValue{Value: code})
	if err != nil {
		return nil, err
	}

	if resp.Suspend != "正常" {
		return nil, fmt.Errorf("invalid stock, suspend: %s", resp.Suspend)
	}
	return &Stock{Code: resp.Code, Name: resp.Name}, nil
}
