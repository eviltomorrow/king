package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/sqlutil"
)

func TransactionRecordWithSelectManyByOpenId(ctx context.Context, exec mysql.Exec, openId string) ([]*TransactionRecord, error) {
	records := make([]*TransactionRecord, 0, 16)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			record := &TransactionRecord{}
			if err := rows.Scan(
				&record.Id,
				&record.UserId,
				&record.FundNo,
				&record.Action,
				&record.AssetsType,
				&record.AssetsCode,
				&record.AssetsName,
				&record.ClosePrice,
				&record.Volumn,
				&record.Datetime,
				&record.Status,
				&record.OpenId,
				&record.CreateTimestamp,
				&record.ModifyTimestamp,
			); err != nil {
				return err
			}
			records = append(records, record)
		}
		return nil
	}
	if err := sqlutil.TableWithSelectMany(ctx, exec, TableTransactionRecordName, TransactionRecordFields, map[string]interface{}{FieldTransactionRecordOpenId: openId}, nil, scan); err != nil {
		return nil, err
	}
	return records, nil
}

func TransactionRecordWithSelectManyByUserId(ctx context.Context, exec mysql.Exec, userId string) ([]*TransactionRecord, error) {
	records := make([]*TransactionRecord, 0, 16)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			record := &TransactionRecord{}
			if err := rows.Scan(
				&record.Id,
				&record.UserId,
				&record.FundNo,
				&record.Action,
				&record.AssetsType,
				&record.AssetsCode,
				&record.AssetsName,
				&record.ClosePrice,
				&record.Volumn,
				&record.Datetime,
				&record.Status,
				&record.OpenId,
				&record.CreateTimestamp,
				&record.ModifyTimestamp,
			); err != nil {
				return err
			}
			records = append(records, record)
		}
		return nil
	}
	if err := sqlutil.TableWithSelectMany(ctx, exec, TableTransactionRecordName, TransactionRecordFields, map[string]interface{}{FieldTransactionRecordUserId: userId}, nil, scan); err != nil {
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
				&record.UserId,
				&record.FundNo,
				&record.Action,
				&record.AssetsType,
				&record.AssetsCode,
				&record.AssetsName,
				&record.ClosePrice,
				&record.Volumn,
				&record.Datetime,
				&record.Status,
				&record.OpenId,
				&record.CreateTimestamp,
				&record.ModifyTimestamp,
			); err != nil {
				return err
			}
			records = append(records, record)
		}
		return nil
	}
	if err := sqlutil.TableWithSelectMany(ctx, exec, TableTransactionRecordName, TransactionRecordFields, map[string]interface{}{FieldTransactionRecordFundNo: fundNo}, nil, scan); err != nil {
		return nil, err
	}
	return records, nil
}

func TransactionRecordWithInsertOne(ctx context.Context, exec mysql.Exec, record *TransactionRecord) (int64, error) {
	if record == nil {
		return 0, fmt.Errorf("record is nil")
	}

	value := map[string]interface{}{
		FieldTransactionRecordId:         record.Id,
		FieldTransactionRecordUserId:     record.UserId,
		FieldTransactionRecordFundNo:     record.FundNo,
		FieldTransactionRecordAction:     record.Action,
		FieldTransactionRecordAssetsType: record.AssetsType,
		FieldTransactionRecordAssetsCode: record.AssetsCode,
		FieldTransactionRecordAssetsName: record.AssetsName,
		FieldTransactionRecordClosePrice: record.ClosePrice,
		FieldTransactionRecordVolume:     record.Volumn,
		FieldTransactionRecordDatetime:   record.Datetime,
		FieldTransactionRecordStatus:     record.Status,
		FieldAssetsOpenId:                record.OpenId,
	}
	return sqlutil.TableWithInsertOne(ctx, exec, TableTransactionRecordName, value)
}

type TransactionRecord struct {
	Id              string       `json:"id"`
	UserId          string       `json:"user_id"`
	FundNo          string       `json:"fund_no"`
	Action          int8         `json:"action"`
	AssetsType      int8         `json:"assets_type"`
	AssetsCode      string       `json:"assets_code"`
	AssetsName      string       `json:"assets_name"`
	ClosePrice      float64      `json:"close_price"`
	Volumn          int64        `json:"volume"`
	Datetime        time.Time    `json:"datetime"`
	Status          int8         `json:"status"`
	OpenId          string       `json:"open_id"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

const (
	TableTransactionRecordName = "TransactionRecord"

	FieldTransactionRecordId              = "id"
	FieldTransactionRecordUserId          = "user_id"
	FieldTransactionRecordFundNo          = "fund_no"
	FieldTransactionRecordAction          = "action"
	FieldTransactionRecordAssetsType      = "assets_type"
	FieldTransactionRecordAssetsCode      = "assets_code"
	FieldTransactionRecordAssetsName      = "assets_name"
	FieldTransactionRecordClosePrice      = "close_price"
	FieldTransactionRecordVolume          = "volume"
	FieldTransactionRecordDatetime        = "datetime"
	FieldTransactionRecordStatus          = "status"
	FieldTransactionRecordOpenId          = "open_id"
	FieldTransactionRecordCreateTimestamp = "create_timestamp"
	FieldTransactionRecordModifyTimestamp = "modify_timestamp"
)

var TransactionRecordFields = []string{
	FieldTransactionRecordId,
	FieldTransactionRecordUserId,
	FieldTransactionRecordFundNo,
	FieldTransactionRecordAction,
	FieldTransactionRecordAssetsType,
	FieldTransactionRecordAssetsCode,
	FieldTransactionRecordClosePrice,
	FieldTransactionRecordVolume,
	FieldTransactionRecordDatetime,
	FieldTransactionRecordStatus,
	FieldTransactionRecordOpenId,
	FieldTransactionRecordCreateTimestamp,
	FieldTransactionRecordModifyTimestamp,
}
