package db

import (
	"database/sql"
	"time"
)

type Fund struct {
	FundNo          string       `json:"fund_no"`
	UserId          string       `json:"user_id"`
	OpeningCash     float64      `json:"opening_cash"`
	EndCash         float64      `json:"end_cash"`
	Status          int8         `json:"status"`
	InitDatetime    time.Time    `json:"init_datetime"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

const (
	FieldFundId              = "id"
	FieldFundFundNo          = "fund_no"
	FieldFundUserId          = "user_id"
	FieldFundOpeningCash     = "opening_cash"
	FieldFundEndCash         = "end_cash"
	FieldFundStatus          = "status"
	FieldFundInitDatetime    = "init_datetime"
	FieldFundCreateTimestamp = "create_timestamp"
	FieldFundModifyTimestamp = "modify_timestamp"
)

var FundFields = []string{
	FieldFundId,
	FieldFundFundNo,
	FieldFundUserId,
	FieldFundOpeningCash,
	FieldFundEndCash,
	FieldFundStatus,
	FieldFundInitDatetime,
}
