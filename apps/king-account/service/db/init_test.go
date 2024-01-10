package db

import (
	"log"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func init() {
	mysql.DSN = "admin:admin123@tcp(mysql:3306)/king_account?charset=utf8mb4&parseTime=true&loc=Local"

	if err := mysql.Connect(); err != nil {
		log.Fatal(err)
	}
}
