package data

import (
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-storage/service/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/model"
)

var timeout = 30 * time.Second

func StoreStock(data []*model.Stock) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	tx, err := mysql.DB.Begin()
	if err != nil {
		return 0, err
	}
	affected, err := db.StockWithInsertOrUpdateMany(tx, data, timeout)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, nil
	}
	return affected, nil
}

func StoreQuote(data []*model.Quote, mode string, date time.Time) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if mode != db.Day && mode != db.Week {
		return 0, fmt.Errorf("invalid mode: %v", mode)
	}

	codes := make([]string, 0, len(data))
	for _, d := range data {
		codes = append(codes, d.Code)
	}

	tx, err := mysql.DB.Begin()
	if err != nil {
		return 0, err
	}

	d := date.Format(time.DateOnly)
	if _, err := db.QuoteWithDeleteManyByCodesAndDate(tx, mode, codes, d, timeout); err != nil {
		tx.Rollback()
		return 0, err
	}

	affected, err := db.QuoteWithInsertMany(tx, mode, data, timeout)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	return affected, nil
}
