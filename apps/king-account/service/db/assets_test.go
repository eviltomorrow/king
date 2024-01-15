package db

import (
	"github.com/eviltomorrow/king/lib/db/mysql"
)

var assets1 = &Assets{
	FundNo: "",
	UserId: "2",
}

func TruncateAssets() error {
	_, err := mysql.DB.Exec("truncate table assets")
	return err
}

func InitAssets() {
	TruncateAssets()
}
