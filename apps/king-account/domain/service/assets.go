package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-account/domain/db"
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
	OpenId           string          `json:"open_id"`
	FirstBuyDatetime time.Time       `json:"first_buy_datetime"`
}

func AssetsWithFindManyByFundNo(ctx context.Context, fundNo string) ([]*Assets, error) {
	if fundNo == "" {
		return nil, fmt.Errorf("fund_no is nil")
	}
	assets, err := db.AssetsWithSelectManyByFundNo(ctx, mysql.DB, fundNo)
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
			OpenId:           asset.OpenId,
			FirstBuyDatetime: asset.FirstBuyDatetime,
		})
	}
	return data, nil
}

func AssetsWithFindManyByUserId(ctx context.Context, userId string) ([]*Assets, error) {
	if userId == "" {
		return nil, fmt.Errorf("userId is nil")
	}
	assets, err := db.AssetsWithSelectManyByUserId(ctx, mysql.DB, userId)
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
			OpenId:           asset.OpenId,
			FirstBuyDatetime: asset.FirstBuyDatetime,
		})
	}
	return data, nil
}

func AssetsWithFindOneByUserIdFundNoCode(ctx context.Context, userId, fundNo, code string) (*Assets, error) {
	if userId == "" || fundNo == "" || code == "" {
		return nil, fmt.Errorf("userId/fundNo/code is nil")
	}

	assets, err := db.AssetsWithSelectOneByUserIdFundNoCode(ctx, mysql.DB, userId, fundNo, code)
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
		OpenId:           assets.OpenId,
		FirstBuyDatetime: assets.FirstBuyDatetime,
	}
	return data, nil
}
