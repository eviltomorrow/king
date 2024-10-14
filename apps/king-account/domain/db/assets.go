package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/sqlutil"
)

func AssetsWithSelectOneByUserIdFundNoCode(ctx context.Context, exec mysql.Exec, userId, fundNo, code string) (*Assets, error) {
	assets := Assets{}
	scan := func(row *sql.Row) error {
		return row.Scan(
			&assets.UserId,
			&assets.FundNo,
			&assets.Type,
			&assets.CashPosition,
			&assets.Code,
			&assets.Name,
			&assets.OpenInterest,
			&assets.OpenId,
			&assets.FirstBuyDatetime,
			&assets.CreateTimestamp,
			&assets.ModifyTimestamp,
		)
	}
	if err := sqlutil.TableWithSelectOne(ctx, exec, TableFundName, FundFields, map[string]interface{}{FieldAssetsFundNo: fundNo, FieldAssetsUserId: userId, FieldAssetsCode: code}, scan); err != nil {
		return nil, err
	}
	return &assets, nil
}

func AssetsWithSelectManyByUserId(ctx context.Context, exec mysql.Exec, userId string) ([]*Assets, error) {
	data := make([]*Assets, 0, 16)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			assets := Assets{}
			if err := rows.Scan(
				&assets.UserId,
				&assets.FundNo,
				&assets.Type,
				&assets.CashPosition,
				&assets.Code,
				&assets.Name,
				&assets.OpenInterest,
				&assets.OpenId,
				&assets.FirstBuyDatetime,
				&assets.CreateTimestamp,
				&assets.ModifyTimestamp,
			); err != nil {
				return err
			}
			data = append(data, &assets)
		}
		return nil
	}
	if err := sqlutil.TableWithSelectMany(ctx, exec, TableAssetsName, AssetsFields, map[string]interface{}{FieldAssetsUserId: userId}, nil, scan); err != nil {
		return nil, err
	}
	return data, nil
}

func AssetsWithSelectManyByFundNo(ctx context.Context, exec mysql.Exec, fundNo string) ([]*Assets, error) {
	data := make([]*Assets, 0, 16)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			assets := Assets{}
			if err := rows.Scan(
				&assets.UserId,
				&assets.FundNo,
				&assets.Type,
				&assets.CashPosition,
				&assets.Code,
				&assets.Name,
				&assets.OpenInterest,
				&assets.OpenId,
				&assets.FirstBuyDatetime,
				&assets.CreateTimestamp,
				&assets.ModifyTimestamp,
			); err != nil {
				return err
			}
			data = append(data, &assets)
		}
		return nil
	}
	if err := sqlutil.TableWithSelectMany(ctx, exec, TableAssetsName, AssetsFields, map[string]interface{}{FieldAssetsFundNo: fundNo}, nil, scan); err != nil {
		return nil, err
	}
	return data, nil
}

func AssetsWithUpdateOneByUserIdFundNoCode(ctx context.Context, exec mysql.Exec, assets *Assets, userId, fundNo, code string) (int64, error) {
	if assets == nil {
		return 0, fmt.Errorf("assets is nil")
	}

	value := map[string]interface{}{
		FieldAssetsName:         assets.Name,
		FieldAssetsCashPosition: assets.CashPosition,
		FieldAssetsOpenInterest: assets.OpenInterest,
	}
	return sqlutil.TableWithUpdate(ctx, exec, TableFundName, value, map[string]interface{}{FieldAssetsFundNo: fundNo, FieldAssetsUserId: userId, FieldAssetsCode: code})
}

func AssetsWithDeleteOneByUserIdFundNoCode(ctx context.Context, exec mysql.Exec, assets *Assets, userId, fundNo, code string) (int64, error) {
	if assets == nil {
		return 0, fmt.Errorf("assets is nil")
	}

	return sqlutil.TableWithDelete(ctx, exec, TableFundName, map[string]interface{}{FieldAssetsFundNo: fundNo, FieldAssetsUserId: userId, FieldAssetsCode: code})
}

func AssetsWithInsertOne(ctx context.Context, exec mysql.Exec, assets *Assets) (int64, error) {
	if assets == nil {
		return 0, fmt.Errorf("assets is nil")
	}

	value := map[string]interface{}{
		FieldAssetsUserId:           assets.UserId,
		FieldAssetsFundNo:           assets.FundNo,
		FieldAssetsType:             assets.Type,
		FieldAssetsCashPosition:     assets.CashPosition,
		FieldAssetsCode:             assets.Code,
		FieldAssetsName:             assets.Name,
		FieldAssetsOpenInterest:     assets.OpenInterest,
		FieldAssetsOpenId:           assets.OpenId,
		FieldAssetsFirstBuyDatetime: assets.FirstBuyDatetime,
	}
	return sqlutil.TableWithInsertOne(ctx, exec, TableAssetsName, value)
}

type Assets struct {
	UserId           string       `json:"user_id"`
	FundNo           string       `json:"fund_no"`
	Type             int8         `json:"type"`
	CashPosition     float64      `json:"cash_position"`
	Code             string       `json:"code"`
	Name             string       `json:"name"`
	OpenInterest     int64        `json:"open_interest"`
	OpenId           string       `json:"open_id"`
	FirstBuyDatetime time.Time    `json:"first_buy_datetime"`
	CreateTimestamp  time.Time    `json:"create_timestamp"`
	ModifyTimestamp  sql.NullTime `json:"modify_timestamp"`
}

const (
	TableAssetsName = "assets"

	FieldAssetsUserId           = "user_id"
	FieldAssetsFundNo           = "fund_no"
	FieldAssetsType             = "type"
	FieldAssetsCashPosition     = "cash_position"
	FieldAssetsCode             = "code"
	FieldAssetsName             = "name"
	FieldAssetsOpenInterest     = "open_interest"
	FieldAssetsOpenId           = "open_id"
	FieldAssetsFirstBuyDatetime = "first_buy_datetime"
	FieldAssetsCreateTimestamp  = "create_timestamp"
	FieldAssetsModifyTimestamp  = "modify_timestamp"
)

var AssetsFields = []string{
	FieldAssetsUserId,
	FieldAssetsFundNo,
	FieldAssetsType,
	FieldAssetsCashPosition,
	FieldAssetsCode,
	FieldAssetsName,
	FieldAssetsOpenInterest,
	FieldAssetsOpenId,
	FieldAssetsFirstBuyDatetime,
	FieldAssetsCreateTimestamp,
	FieldAssetsModifyTimestamp,
}
