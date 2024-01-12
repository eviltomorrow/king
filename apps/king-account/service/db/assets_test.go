package db

import (
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func TruncateAssets() error {
	_, err := mysql.DB.Exec("truncate table assets")
	return err
}

func InitAssets() {
	TruncateAccount()
}

var assets1 = &Assets{
	Id:               "",
	FundNo:           "",
	UserId:           "",
	Type:             1,
	CashPosition:     10,
	Code:             "",
	OpenInterest:     20,
	FirstBuyDatetime: time.Date(2024, time.January, 10, 10, 0o0, 0o0, 0o0, time.Local),
}
