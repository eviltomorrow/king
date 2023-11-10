package storage

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

func StoreMetadata(date time.Time, metadata chan *model.Metadata) (int64, int64, int64, error) {
	var (
		affectedS, affectedD, affectedW int64
		i, size                         = 0, 50
	)

	data := make([]*model.Metadata, 0, size)
	for md := range metadata {
		if md == nil {
			continue
		}
		if i >= size {
			s, d, w, err := storeMetadata(date, data)
			if err != nil {
				return 0, 0, 0, err
			}
			affectedS, affectedD, affectedW = affectedS+s, affectedD+d, affectedW+w
			i = 0
			data = data[:0]
		} else {
			data = append(data, md)
		}
	}
	if len(data) > 0 {
		s, d, w, err := storeMetadata(date, data)
		if err != nil {
			return 0, 0, 0, err
		}
		affectedS, affectedD, affectedW = affectedS+s, affectedD+d, affectedW+w
	}

	return affectedS, affectedD, affectedW, nil
}

func storeMetadata(date time.Time, metadata []*model.Metadata) (int64, int64, int64, error) {
	var (
		stocks = make([]*model.Stock, 0, len(metadata))
		days   = make([]*model.Quote, 0, len(metadata))

		affectedS, affectedD, affectedW int64
	)

	for _, md := range metadata {
		stocks = append(stocks, &model.Stock{
			Code:            md.Code,
			Name:            md.Name,
			Suspend:         md.Suspend,
			CreateTimestamp: time.Now(),
		})

		day, err := AssembleQuoteDay(md, date)
		if err != nil {
			return 0, 0, 0, err
		}
		days = append(days, day)
	}

	if len(stocks) != 0 {
		affected, err := StoreStock(stocks)
		if err != nil {
			return 0, 0, 0, err
		}
		affectedS += affected

		affected, err = StoreQuote(days, db.Day, date)
		if err != nil {
			return 0, 0, 0, err
		}
		affectedD += affected
	}

	if date.Weekday() == time.Friday {
		weeks := make([]*model.Quote, 0, len(stocks))
		for _, stock := range stocks {
			week, err := AssembleQuoteWeek(stock.Code, date)
			if err != nil && err != ErrNoData {
				return affectedS, affectedD, affectedW, err
			}
			if err == ErrNoData {
				continue
			}
			if week != nil {
				weeks = append(weeks, week)
			}
		}
		affected, err := StoreQuote(weeks, db.Week, date)
		if err != nil {
			return affectedS, affectedD, affectedW, err
		}
		affectedW += affected
	}
	return affectedS, affectedD, affectedW, nil
}
