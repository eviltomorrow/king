package impl

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-account/domain/service"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-account"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type AssetsServer struct {
	pb.UnimplementedAssetsServer
}

func (s *AssetsServer) ListByUserId(ctx context.Context, req *wrapperspb.StringValue) (*pb.ItemListResp, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	assets, err := service.AssetsWithFindManyByUserId(ctx, req.Value)
	if err != nil {
		return nil, err
	}

	count := len(assets)
	data := make([]*pb.Item, 0, count)
	for _, asset := range assets {
		assetType := pb.Item_STOCK

		data = append(data, &pb.Item{
			FundNo:           asset.FundNo,
			UserId:           asset.UserId,
			Type:             assetType,
			CashPosition:     asset.CashPosition.InexactFloat64(),
			Code:             asset.Code,
			OpenInterest:     asset.OpenInterest,
			FirstBuyDatetime: asset.FirstBuyDatetime.Unix(),
		})
	}
	return &pb.ItemListResp{TotalCount: int64(count), Items: data}, nil
}
