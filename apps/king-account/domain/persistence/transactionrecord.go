package persistence

import (
	"database/sql"
	"time"
)

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

	FieldTransactionRecordId               = "id"
	FieldTransactionRecordUsername         = "username"
	FieldTransactionRecordPassword         = "password"
	FieldTransactionRecordNickName         = "nick_name"
	FieldTransactionRecordPhone            = "phone"
	FieldTransactionRecordEmail            = "email"
	FieldTransactionRecordStatus           = "status"
	FieldTransactionRecordRegisterDatetime = "register_datetime"
	FieldTransactionRecordCreateTimestamp  = "create_timestamp"
	FieldTransactionRecordModifyTimestamp  = "modify_timestamp"
)
