package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-account/domain/persistence"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/shopspring/decimal"
)

type Assets struct {
	Id               string          `json:"id"`
	FundNo           string          `json:"fund_no"`
	UserId           string          `json:"user_id"`
	Type             int8            `json:"type"`
	CashPosition     decimal.Decimal `json:"cash_position"`
	Code             string          `json:"code"`
	OpenInterest     int64           `json:"open_interest"`
	FirstBuyDatetime time.Time       `json:"first_buy_datetime"`
}

func AssetsWithFindByUserId(ctx context.Context, userId string) ([]*Assets, error) {
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
			Id:               asset.Id,
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

func AssetsWithSelectOneByUserIdFundNoCode(ctx context.Context, userId, fundNo, code string) (*Assets, error) {
	if userId == "" || fundNo == "" || code == "" {
		return nil, fmt.Errorf("userId/fundNo/code is nil")
	}

	assets, err := persistence.AssetsWithSelectOneByUserIdFundNoCode(ctx, mysql.DB, userId, fundNo, code)
	if err != nil {
		return nil, err
	}

	d := decimal.NewFromFloat(assets.CashPosition)
	data := &Assets{
		Id:               assets.Id,
		FundNo:           assets.FundNo,
		UserId:           assets.UserId,
		Type:             assets.Type,
		CashPosition:     d,
		Code:             assets.Code,
		OpenInterest:     assets.OpenInterest,
		FirstBuyDatetime: assets.FirstBuyDatetime,
	}
	return data, nil
}

func AssetsWithBuy(ctx context.Context, assets *Assets) error {
	if assets == nil {
		return fmt.Errorf("assets is nil")
	}

	currentAssets, err := AssetsWithSelectOneByUserIdFundNoCode(ctx, assets.UserId, assets.FundNo, assets.Code)
	if err == sql.ErrNoRows {
		_, err = persistence.AssetsWithInsertOne(ctx, mysql.DB, &persistence.Assets{
			UserId:           assets.UserId,
			FundNo:           assets.FundNo,
			Type:             assets.Type,
			CashPosition:     assets.CashPosition.InexactFloat64(),
			Code:             assets.Code,
			OpenInterest:     assets.OpenInterest,
			FirstBuyDatetime: time.Now(),
		})
		return err
	}
	if err != nil {
		return err
	}

	newAssets := &persistence.Assets{
		UserId: currentAssets.UserId,
		FundNo: currentAssets.FundNo,
		Type:   currentAssets.Type,
		Code:   currentAssets.Code,
	}
	newAssets.OpenInterest += assets.OpenInterest
	newAssets.CashPosition = currentAssets.CashPosition.Add(assets.CashPosition).InexactFloat64()

	_, err = persistence.AssetsWithUpdateOneByUserIdFundNoCode(ctx, mysql.DB, newAssets, currentAssets.UserId, currentAssets.FundNo, currentAssets.Code)
	return err
}

func AssetsWithSell(ctx context.Context, assets *Assets) error {
	if assets == nil {
		return fmt.Errorf("assets is nil")
	}

	currentAssets, err := AssetsWithSelectOneByUserIdFundNoCode(ctx, assets.UserId, assets.FundNo, assets.Code)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no assets")
	}
	if err != nil {
		return err
	}

	newAssets := &persistence.Assets{
		UserId: currentAssets.UserId,
		FundNo: currentAssets.FundNo,
		Type:   currentAssets.Type,
		Code:   currentAssets.Code,
	}
	newAssets.OpenInterest -= assets.OpenInterest
	newAssets.CashPosition = currentAssets.CashPosition.Sub(assets.CashPosition).InexactFloat64()

	if newAssets.OpenInterest < 0 {
		return fmt.Errorf("OpenInterest has oversold")
	}
	if newAssets.CashPosition < 0 {
		return fmt.Errorf("CashPosition has oversold")
	}

	if newAssets.OpenInterest > 0 || newAssets.CashPosition > 0 {
		_, err = persistence.AssetsWithUpdateOneByUserIdFundNoCode(ctx, mysql.DB, newAssets, currentAssets.UserId, currentAssets.FundNo, currentAssets.Code)
		return err
	}

	if newAssets.OpenInterest == 0 || newAssets.CashPosition == 0 {
		_, err = persistence.AssetsWithDeleteOneByUserIdFundNoCode(ctx, mysql.DB, newAssets, currentAssets.UserId, currentAssets.FundNo, currentAssets.Code)
		return err
	}

	return nil
}
