package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/orm"
)

func AssetsWithSelectManyByUserId(ctx context.Context, exec mysql.Exec, userId string) ([]*Assets, error) {
	data := make([]*Assets, 0, 16)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			assets := Assets{}
			if err := rows.Scan(
				&assets.Id,
				&assets.FundNo,
				&assets.UserId,
				&assets.Type,
				&assets.CashPosition,
				&assets.Code,
				&assets.OpenInterest,
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
	if err := orm.TableWithSelectMany(ctx, exec, TableAssetsName, AssetsFields, map[string]interface{}{FieldAssetsUserId: userId}, nil, scan); err != nil {
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
				&assets.Id,
				&assets.FundNo,
				&assets.UserId,
				&assets.Type,
				&assets.CashPosition,
				&assets.Code,
				&assets.OpenInterest,
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
	if err := orm.TableWithSelectMany(ctx, exec, TableAssetsName, AssetsFields, map[string]interface{}{FieldAssetsFundNo: fundNo}, nil, scan); err != nil {
		return nil, err
	}
	return data, nil
}

func AssetsWithUpdateOne(ctx context.Context, exec mysql.Exec, assets *Assets, id string) (int64, error) {
	if assets == nil {
		return 0, fmt.Errorf("invalid parameter, value is nil")
	}

	value := map[string]interface{}{
		FieldAssetsCashPosition: assets.CashPosition,
		FieldAssetsOpenInterest: assets.OpenInterest,
	}
	return orm.TableWithUpdate(ctx, exec, TableFundName, value, map[string]interface{}{FieldAssetsId: id})
}

func AssetsWithUpdateOneByFundNoUserIdCode(ctx context.Context, exec mysql.Exec, assets *Assets, fundNo, userId, code string) (int64, error) {
	if assets == nil {
		return 0, fmt.Errorf("invalid parameter, assets is nil")
	}

	value := map[string]interface{}{
		FieldAssetsCashPosition: assets.CashPosition,
		FieldAssetsOpenInterest: assets.OpenInterest,
	}
	return orm.TableWithUpdate(ctx, exec, TableFundName, value, map[string]interface{}{FieldAssetsFundNo: fundNo, FieldAssetsUserId: userId, FieldAssetsCode: code})
}

func AssetsWithDeleteOne(ctx context.Context, exec mysql.Exec, id string) (int64, error) {
	if id == "" {
		return 0, fmt.Errorf("invalid parameter, id is nil")
	}
	return orm.TableWithDelete(ctx, exec, TableAssetsName, map[string]interface{}{FieldAssetsId: id})
}

func AssetsWithInsertOne(ctx context.Context, exec mysql.Exec, assets *Assets) (int64, error) {
	if assets == nil {
		return 0, fmt.Errorf("invalid parameter, assets is nil")
	}

	value := map[string]interface{}{
		FieldAssetsId:               assets.Id,
		FieldAssetsFundNo:           assets.FundNo,
		FieldAssetsUserId:           assets.UserId,
		FieldAssetsType:             assets.Type,
		FieldAssetsCashPosition:     assets.CashPosition,
		FieldAssetsCode:             assets.Code,
		FieldAssetsOpenInterest:     assets.OpenInterest,
		FieldAssetsFirstBuyDatetime: assets.FirstBuyDatetime,
	}
	return orm.TableWithInsert(ctx, exec, TableAssetsName, value)
}

type Assets struct {
	Id               string       `json:"id"`
	FundNo           string       `json:"fund_no"`
	UserId           string       `json:"user_id"`
	Type             int8         `json:"type"`
	CashPosition     float64      `json:"cash_position"`
	Code             string       `json:"code"`
	OpenInterest     float64      `json:"open_interest"`
	FirstBuyDatetime time.Time    `json:"first_buy_datetime"`
	CreateTimestamp  time.Time    `json:"create_timestamp"`
	ModifyTimestamp  sql.NullTime `json:"modify_timestamp"`
}

const (
	TableAssetsName = "assets"

	FieldAssetsId               = "id"
	FieldAssetsFundNo           = "fund_no"
	FieldAssetsUserId           = "user_id"
	FieldAssetsType             = "type"
	FieldAssetsCashPosition     = "cash_position"
	FieldAssetsCode             = "code"
	FieldAssetsOpenInterest     = "open_interest"
	FieldAssetsFirstBuyDatetime = "first_buy_datetime"
	FieldAssetsCreateTimestamp  = "create_timestamp"
	FieldAssetsModifyTimestamp  = "modify_timestamp"
)

var AssetsFields = []string{
	FieldAssetsId,
	FieldAssetsFundNo,
	FieldAssetsUserId,
	FieldAssetsType,
	FieldAssetsCashPosition,
	FieldAssetsCode,
	FieldAssetsOpenInterest,
	FieldAssetsFirstBuyDatetime,
	FieldAssetsCreateTimestamp,
	FieldAssetsModifyTimestamp,
}
