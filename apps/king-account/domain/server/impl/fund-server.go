package impl

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-account/domain/service"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-account"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type FundServer struct {
	pb.UnimplementedFundServer
}

// ReblanceByUserId(context.Context, *wrapperspb.StringValue) (*emptypb.Empty, error)
func (s *FundServer) InitAccount(ctx context.Context, req *pb.InitAccountReq) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	if err := service.FundWithInitAccount(ctx, req.AliasName, req.UserId, decimal.NewFromFloat(req.OpeningCash)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *FundServer) ListByUserId(ctx context.Context, req *wrapperspb.StringValue) (*pb.AccountListResp, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	accounts, err := service.FundWithListByUserId(ctx, req.Value)
	if err != nil {
		return nil, err
	}

	resp := &pb.AccountListResp{
		Accounts:   make([]*pb.Account, 0, len(accounts)),
		TotalCount: int64(len(accounts)),
	}
	for _, account := range accounts {
		resp.Accounts = append(resp.Accounts, &pb.Account{
			AliasName:        account.AliasName,
			UserId:           account.UserId,
			FundNo:           account.FundNo,
			OpeningCash:      account.OpeningCash.InexactFloat64(),
			EndCash:          account.EndCash.InexactFloat64(),
			YesterdayEndCash: account.YesterdayEndCash.InexactFloat64(),
			Status:           int32(account.Status),
			InitDatetime:     account.InitDatetime.Unix(),
		})
	}
	return resp, nil
}

func (s *AssetsServer) ModifyAliasName(ctx context.Context, req *pb.ModifyAliasNameReq) (*emptypb.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	if err := service.FundWithModifyAliasName(ctx, req.AliasName, req.FundNo); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
