package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/orm"
)

func FundWithCountByUserId(ctx context.Context, exec mysql.Exec, userId string) (int64, error) {
	if userId == "" {
		return 0, fmt.Errorf("user_id is nil")
	}

	return orm.TableWithCount(ctx, exec, TableFundName, map[string]interface{}{FieldFundUserId: userId})
}

func FundWithSelectOne(ctx context.Context, exec mysql.Exec, fundNo string) (*Fund, error) {
	fund := Fund{}
	scan := func(row *sql.Row) error {
		return row.Scan(
			&fund.AliasName,
			&fund.UserId,
			&fund.FundNo,
			&fund.OpeningCash,
			&fund.EndCash,
			&fund.YesterdayEndCash,
			&fund.Status,
			&fund.InitDatetime,
			&fund.Version,
			&fund.CreateTimestamp,
			&fund.ModifyTimestamp,
		)
	}
	if err := orm.TableWithSelectOne(ctx, exec, TableFundName, FundFields, map[string]interface{}{FieldFundFundNo: fundNo}, scan); err != nil {
		return nil, err
	}
	return &fund, nil
}

func FundWithSelectManyByUserId(ctx context.Context, exec mysql.Exec, userId string) ([]*Fund, error) {
	funds := make([]*Fund, 0, 8)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			fund := Fund{}
			if err := rows.Scan(
				&fund.AliasName,
				&fund.UserId,
				&fund.FundNo,
				&fund.OpeningCash,
				&fund.EndCash,
				&fund.YesterdayEndCash,
				&fund.Status,
				&fund.InitDatetime,
				&fund.Version,
				&fund.CreateTimestamp,
				&fund.ModifyTimestamp,
			); err != nil {
				return err
			}
			funds = append(funds, &fund)
		}
		return nil
	}
	if err := orm.TableWithSelectMany(ctx, exec, TableFundName, FundFields, map[string]interface{}{FieldFundUserId: userId}, nil, scan); err != nil {
		return nil, err
	}
	return funds, nil
}

func FundWithUpdateAliasName(ctx context.Context, exec mysql.Exec, aliasName, fundNo string) (int64, error) {
	if aliasName == "" {
		return 0, fmt.Errorf("alias_name is nil")
	}
	if fundNo == "" {
		return 0, fmt.Errorf("fund_no is nil")
	}

	value := map[string]interface{}{
		FieldFundAliasName: aliasName,
	}
	return orm.TableWithUpdate(ctx, exec, TableFundName, value, map[string]interface{}{FieldFundFundNo: fundNo})
}

func FundWithUpdateOne(ctx context.Context, exec mysql.Exec, fund *Fund, fundNo string, version int64) (int64, error) {
	if fund == nil {
		return 0, fmt.Errorf("fund is nil")
	}
	if fundNo == "" {
		return 0, fmt.Errorf("fund_no is nil")
	}

	value := map[string]interface{}{
		FieldFundAliasName:        fund.AliasName,
		FieldFundOpeningCash:      fund.OpeningCash,
		FieldFundEndCash:          fund.EndCash,
		FieldFundYesterdayEndCash: fund.YesterdayEndCash,
		FieldFundStatus:           fund.Status,
		FieldFundVersion:          fund.Version + 1,
	}
	return orm.TableWithUpdate(ctx, exec, TableFundName, value, map[string]interface{}{FieldFundFundNo: fundNo, FieldFundVersion: version})
}

func FundWithDeleteOne(ctx context.Context, exec mysql.Exec, fundNo string) (int64, error) {
	if fundNo == "" {
		return 0, fmt.Errorf("fund_no is nil")
	}
	return orm.TableWithDelete(ctx, exec, TableFundName, map[string]interface{}{FieldFundFundNo: fundNo})
}

func FundWithInsertOne(ctx context.Context, exec mysql.Exec, fund *Fund) (int64, error) {
	if fund == nil {
		return 0, fmt.Errorf("fund is nil")
	}

	value := map[string]interface{}{
		FieldFundAliasName:        fund.AliasName,
		FieldFundUserId:           fund.UserId,
		FieldFundFundNo:           fund.FundNo,
		FieldFundOpeningCash:      fund.OpeningCash,
		FieldFundEndCash:          fund.EndCash,
		FieldFundYesterdayEndCash: fund.YesterdayEndCash,
		FieldFundStatus:           fund.Status,
		FieldFundInitDatetime:     fund.InitDatetime,
		FieldFundVersion:          0,
	}
	return orm.TableWithInsertOne(ctx, exec, TableFundName, value)
}

type Fund struct {
	AliasName        string       `json:"alias_name"`
	UserId           string       `json:"user_id"`
	FundNo           string       `json:"fund_no"`
	OpeningCash      float64      `json:"opening_cash"`
	EndCash          float64      `json:"end_cash"`
	YesterdayEndCash float64      `json:"yesterday_end_cash"`
	Status           int8         `json:"status"`
	InitDatetime     time.Time    `json:"init_datetime"`
	Version          int64        `json:"version"`
	CreateTimestamp  time.Time    `json:"create_timestamp"`
	ModifyTimestamp  sql.NullTime `json:"modify_timestamp"`
}

const (
	TableFundName = "fund"

	FieldFundAliasName        = "alias_name"
	FieldFundUserId           = "user_id"
	FieldFundFundNo           = "fund_no"
	FieldFundOpeningCash      = "opening_cash"
	FieldFundEndCash          = "end_cash"
	FieldFundYesterdayEndCash = "yesterday_end_cash"
	FieldFundStatus           = "status"
	FieldFundInitDatetime     = "init_datetime"
	FieldFundVersion          = "version"
	FieldFundCreateTimestamp  = "create_timestamp"
	FieldFundModifyTimestamp  = "modify_timestamp"
)

var FundFields = []string{
	FieldFundAliasName,
	FieldFundUserId,
	FieldFundFundNo,
	FieldFundOpeningCash,
	FieldFundEndCash,
	FieldFundYesterdayEndCash,
	FieldFundStatus,
	FieldFundInitDatetime,
	FieldFundVersion,
	FieldFundCreateTimestamp,
	FieldFundModifyTimestamp,
}
