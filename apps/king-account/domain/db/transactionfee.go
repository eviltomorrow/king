package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/sqlutil"
)

func TransactionFeeWithSelectManyByRecordId(ctx context.Context, exec mysql.Exec, recordId string) (*TransactionFee, error) {
	fee := TransactionFee{}
	scan := func(row *sql.Row) error {
		return row.Scan(
			&fee.Id,
			&fee.RecordId,
			&fee.Item,
			&fee.Money,
			&fee.CreateTimestamp,
			&fee.ModifyTimestamp,
		)
	}
	if err := sqlutil.TableWithSelectOne(ctx, exec, TableTransactionFeeName, TransactionFeeFields, map[string]interface{}{FieldTransactionFeeRecordId: recordId}, scan); err != nil {
		return nil, err
	}
	return &fee, nil
}

func TransactionFeeWithInsertMany(ctx context.Context, exec mysql.Exec, fees []*TransactionFee) (int64, error) {
	if len(fees) == 0 {
		return 0, fmt.Errorf("fees is nil")
	}

	return 0, nil
}

func TransactionFeeWithInsertOne(ctx context.Context, exec mysql.Exec, fee *TransactionFee) (int64, error) {
	if fee == nil {
		return 0, fmt.Errorf("fee is nil")
	}

	value := map[string]interface{}{
		FieldTransactionFeeId:       fee.Id,
		FieldTransactionFeeRecordId: fee.RecordId,
		FieldTransactionFeeItem:     fee.Item,
		FieldTransactionFeeMoney:    fee.Money,
	}
	return sqlutil.TableWithInsertOne(ctx, exec, TableTransactionFeeName, value)
}

type TransactionFee struct {
	Id              string       `json:"id"`
	RecordId        string       `json:"record_id"`
	Item            string       `json:"item"`
	Money           float64      `json:"money"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

const (
	TableTransactionFeeName = "transaction_fee"

	FieldTransactionFeeId              = "id"
	FieldTransactionFeeRecordId        = "record_id"
	FieldTransactionFeeItem            = "item"
	FieldTransactionFeeMoney           = "money"
	FieldTransactionFeeCreateTimestamp = "create_timestamp"
	FieldTransactionFeeModifyTimestamp = "modify_timestamp"
)

var TransactionFeeFields = []string{
	FieldTransactionFeeId,
	FieldTransactionFeeRecordId,
	FieldTransactionFeeItem,
	FieldTransactionFeeMoney,
	FieldTransactionFeeCreateTimestamp,
	FieldTransactionFeeModifyTimestamp,
}
