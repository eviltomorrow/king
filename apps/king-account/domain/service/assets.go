package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-account/domain/persistence"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/shopspring/decimal"
)

type Assets struct {
	UserId           string          `json:"user_id"`
	FundNo           string          `json:"fund_no"`
	Type             int8            `json:"type"`
	CashPosition     decimal.Decimal `json:"cash_position"`
	Code             string          `json:"code"`
	Name             string          `json:"name"`
	OpenInterest     int64           `json:"open_interest"`
	FirstBuyDatetime time.Time       `json:"first_buy_datetime"`
}

func AssetsWithFindManyByUserId(ctx context.Context, userId string) ([]*Assets, error) {
	if userId == "" {
		return nil, fmt.Errorf("userId is nil")
	}
	assets, err := persistence.AssetsWithSelectManyByUserId(ctx, mysql.DB, userId)
	if err != nil {
		return nil, err
	}

	data := make([]*Assets, 0, len(assets))
	for _, asset := range assets {
		d := decimal.NewFromFloat(asset.CashPosition)

		data = append(data, &Assets{
			FundNo:           asset.FundNo,
			UserId:           asset.UserId,
			Type:             asset.Type,
			CashPosition:     d,
			Code:             asset.Code,
			OpenInterest:     asset.OpenInterest,
			FirstBuyDatetime: asset.FirstBuyDatetime,
		})
	}
	return data, nil
}

func AssetsWithFindOneByUserIdFundNoCode(ctx context.Context, userId, fundNo, code string) (*Assets, error) {
	if userId == "" || fundNo == "" || code == "" {
		return nil, fmt.Errorf("userId/fundNo/code is nil")
	}

	assets, err := persistence.AssetsWithSelectOneByUserIdFundNoCode(ctx, mysql.DB, userId, fundNo, code)
	if err != nil {
		return nil, err
	}

	d := decimal.NewFromFloat(assets.CashPosition)
	data := &Assets{
		UserId:           assets.UserId,
		FundNo:           assets.FundNo,
		Type:             assets.Type,
		CashPosition:     d,
		Code:             assets.Code,
		Name:             assets.Name,
		OpenInterest:     assets.OpenInterest,
		FirstBuyDatetime: assets.FirstBuyDatetime,
	}
	return data, nil
}
