package db

import (
	"database/sql"
	"time"
)

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
