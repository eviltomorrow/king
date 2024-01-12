package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func AssetsWithSelectManyByUserId(exec mysql.Exec, userId string, timeout time.Duration) ([]*Assets, error) {
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	_sql := `select id, fund_no, user_id, type, cash_position, code, open_interest, first_buy_datetime, create_timestamp, modify_timestamp from account limit ?, ?`
	rows, err := exec.QueryContext(ctx, _sql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]*Assets, 0, 8)
	for rows.Next() {
		a := Assets{}
		if err := rows.Scan(
			&a.Id,
			&a.FundNo,
			&a.UserId,
			&a.Type,
			&a.CashPosition,
			&a.Code,
			&a.OpenInterest,
			&a.FirstBuyDatetime,
			&a.CreateTimestamp,
			&a.ModifyTimestamp,
		); err != nil {
			return nil, err
		}
		data = append(data, &a)
	}
	return data, nil
}

func AssetsWithDeleteOne(exec mysql.Exec, id string, timeout time.Duration) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_sql := `delete from assets where id = ?`

	result, err := exec.ExecContext(ctx, _sql, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func AssetsWithInsertOne(exec mysql.Exec, assets *Assets, timeout time.Duration) (int64, error) {
	if assets == nil {
		return 0, fmt.Errorf("invalid parameter, assets is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fields := `(?, ?, ?, ?, ?, ?, ?, ?, now())`
	args := []interface{}{
		assets.Id,
		assets.FundNo,
		assets.UserId,
		assets.Type,
		assets.CashPosition,
		assets.Code,
		assets.OpenInterest,
		assets.FirstBuyDatetime,
	}
	_sql := fmt.Sprintf(`insert into assets (%s) values %s`, strings.Join(AssetsFields, ","), fields)
	result, err := exec.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
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
}
