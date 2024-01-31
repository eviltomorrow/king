package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/orm"
)

func FundWithSelectOne(ctx context.Context, exec mysql.Exec, fundNo string) (*Fund, error) {
	fund := Fund{}
	scan := func(row *sql.Row) error {
		return row.Scan(
			&fund.UserId,
			&fund.FundNo,
			&fund.OpeningCash,
			&fund.EndCash,
			&fund.Status,
			&fund.InitDatetime,
			&fund.CreateTimestamp,
			&fund.ModifyTimestamp,
		)
	}
	if err := orm.TableWithSelectOne(ctx, exec, TableFundName, FundFields, map[string]interface{}{FieldFundFundNo: fundNo}, scan); err != nil {
		return nil, err
	}
	return &fund, nil
}

func FundWithSelectRangeByUserId(ctx context.Context, exec mysql.Exec, userId string, offset, limit int64) ([]*Fund, error) {
	funds := make([]*Fund, 0, limit)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			fund := Fund{}
			if err := rows.Scan(
				&fund.UserId,
				&fund.FundNo,
				&fund.OpeningCash,
				&fund.EndCash,
				&fund.Status,
				&fund.InitDatetime,
				&fund.CreateTimestamp,
				&fund.ModifyTimestamp,
			); err != nil {
				return err
			}
			funds = append(funds, &fund)
		}
		return nil
	}
	if err := orm.TableWithSelectRange(ctx, exec, TableFundName, FundFields, map[string]interface{}{FieldFundUserId: userId}, nil, offset, limit, scan); err != nil {
		return nil, err
	}
	return funds, nil
}

func FundWithUpdateOne(ctx context.Context, exec mysql.Exec, fund *Fund, fundNo string) (int64, error) {
	if fund == nil {
		return 0, fmt.Errorf("invalid parameter, fund is nil")
	}
	if fundNo == "" {
		return 0, fmt.Errorf("invalid parameter, fund_no is nil")
	}

	value := map[string]interface{}{
		FieldFundOpeningCash: fund.OpeningCash,
		FieldFundEndCash:     fund.EndCash,
		FieldFundStatus:      fund.Status,
	}
	return orm.TableWithUpdate(ctx, exec, TableFundName, value, map[string]interface{}{FieldFundFundNo: fundNo})
}

func FundWithDeleteOne(ctx context.Context, exec mysql.Exec, fundNo string) (int64, error) {
	if fundNo == "" {
		return 0, fmt.Errorf("invalid parameter, fund_no is nil")
	}
	return orm.TableWithDelete(ctx, exec, TableFundName, map[string]interface{}{FieldFundFundNo: fundNo})
}

func FundWithInsertOne(ctx context.Context, exec mysql.Exec, fund *Fund) (int64, error) {
	if fund == nil {
		return 0, fmt.Errorf("invalid parameter, fund is nil")
	}

	value := map[string]interface{}{
		FieldFundUserId:       fund.UserId,
		FieldFundFundNo:       fund.FundNo,
		FieldFundOpeningCash:  fund.OpeningCash,
		FieldFundEndCash:      fund.EndCash,
		FieldFundStatus:       fund.Status,
		FieldFundInitDatetime: fund.InitDatetime,
	}
	return orm.TableWithInsert(ctx, exec, TableFundName, value)
}

type Fund struct {
	UserId          string       `json:"user_id"`
	FundNo          string       `json:"fund_no"`
	OpeningCash     float64      `json:"opening_cash"`
	EndCash         float64      `json:"end_cash"`
	Status          int8         `json:"status"`
	InitDatetime    time.Time    `json:"init_datetime"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

const (
	TableFundName = "fund"

	FieldFundUserId          = "user_id"
	FieldFundFundNo          = "fund_no"
	FieldFundOpeningCash     = "opening_cash"
	FieldFundEndCash         = "end_cash"
	FieldFundStatus          = "status"
	FieldFundInitDatetime    = "init_datetime"
	FieldFundCreateTimestamp = "create_timestamp"
	FieldFundModifyTimestamp = "modify_timestamp"
)

var FundFields = []string{
	FieldFundUserId,
	FieldFundFundNo,
	FieldFundOpeningCash,
	FieldFundEndCash,
	FieldFundStatus,
	FieldFundInitDatetime,
	FieldFundCreateTimestamp,
	FieldFundModifyTimestamp,
}
