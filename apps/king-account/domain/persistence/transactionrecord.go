package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/orm"
)

func TransactionRecordWithSelectManyByUserId(ctx context.Context, exec mysql.Exec, userId string) ([]*TransactionRecord, error) {
	records := make([]*TransactionRecord, 0, 16)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			record := &TransactionRecord{}
			if err := rows.Scan(
				&record.Id,
				&record.FundNo,
				&record.UserId,
				&record.Action,
				&record.AssetsType,
				&record.AssetsCode,
				&record.ClosePrice,
				&record.Volumn,
				&record.Datetime,
				&record.CreateTimestamp,
				&record.ModifyTimestamp,
			); err != nil {
				return err
			}
			records = append(records, record)
		}
		return nil
	}
	if err := orm.TableWithSelectMany(ctx, exec, TableTransactionRecordName, TransactionRecordFields, map[string]interface{}{FieldTransactionRecordUserId: userId}, nil, scan); err != nil {
		return nil, err
	}
	return records, nil
}

func TransactionRecordWithSelectManyByFundNo(ctx context.Context, exec mysql.Exec, fundNo string) ([]*TransactionRecord, error) {
	records := make([]*TransactionRecord, 0, 16)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			record := &TransactionRecord{}
			if err := rows.Scan(
				&record.Id,
				&record.FundNo,
				&record.UserId,
				&record.Action,
				&record.AssetsType,
				&record.AssetsCode,
				&record.ClosePrice,
				&record.Volumn,
				&record.Datetime,
				&record.CreateTimestamp,
				&record.ModifyTimestamp,
			); err != nil {
				return err
			}
			records = append(records, record)
		}
		return nil
	}
	if err := orm.TableWithSelectMany(ctx, exec, TableTransactionRecordName, TransactionRecordFields, map[string]interface{}{FieldTransactionRecordFundNo: fundNo}, nil, scan); err != nil {
		return nil, err
	}
	return records, nil
}

func TransactionRecordWithInsertOne(ctx context.Context, exec mysql.Exec, record *TransactionRecord) (int64, error) {
	if record == nil {
		return 0, fmt.Errorf("invalid parameter, record is nil")
	}

	value := map[string]interface{}{
		FieldTransactionRecordId:         record.Id,
		FieldTransactionRecordFundNo:     record.FundNo,
		FieldTransactionRecordUserId:     record.UserId,
		FieldTransactionRecordAction:     record.Action,
		FieldTransactionRecordAssetsType: record.AssetsType,
		FieldTransactionRecordAssetsCode: record.AssetsCode,
		FieldTransactionRecordClosePrice: record.ClosePrice,
		FieldTransactionRecordVolume:     record.Volumn,
		FieldTransactionRecordDatetime:   record.Datetime,
	}
	return orm.TableWithInsert(ctx, exec, TableTransactionRecordName, value)
}

type TransactionRecord struct {
	Id              string       `json:"id"`
	FundNo          string       `json:"fund_no"`
	UserId          string       `json:"user_id"`
	Action          int8         `json:"action"`
	AssetsType      int8         `json:"assets_type"`
	AssetsCode      string       `json:"assets_code"`
	ClosePrice      float64      `json:"close_price"`
	Volumn          int64        `json:"volume"`
	Datetime        time.Time    `json:"datetime"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

const (
	TableTransactionRecordName = "TransactionRecord"

	FieldTransactionRecordId              = "id"
	FieldTransactionRecordFundNo          = "fund_no"
	FieldTransactionRecordUserId          = "user_id"
	FieldTransactionRecordAction          = "action"
	FieldTransactionRecordAssetsType      = "assets_type"
	FieldTransactionRecordAssetsCode      = "assets_code"
	FieldTransactionRecordClosePrice      = "close_price"
	FieldTransactionRecordVolume          = "volume"
	FieldTransactionRecordDatetime        = "datetime"
	FieldTransactionRecordCreateTimestamp = "create_timestamp"
	FieldTransactionRecordModifyTimestamp = "modify_timestamp"
)

var TransactionRecordFields = []string{
	FieldTransactionRecordId,
	FieldTransactionRecordFundNo,
	FieldTransactionRecordUserId,
	FieldTransactionRecordAction,
	FieldTransactionRecordAssetsType,
	FieldTransactionRecordAssetsCode,
	FieldTransactionRecordClosePrice,
	FieldTransactionRecordVolume,
	FieldTransactionRecordDatetime,
	FieldTransactionRecordCreateTimestamp,
	FieldTransactionRecordModifyTimestamp,
}
