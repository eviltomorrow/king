package service

import (
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-storage/domain/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/model"
)

var DBExecTimeout = 30 * time.Second

func StoreMetadata(date time.Time, metadata []*model.Metadata) (int64, int64, error) {
	var (
		affectedStock, affectedQuote int64
		i, size                      = 0, 30
	)

	data := make([]*model.Metadata, 0, size)
	for _, md := range metadata {
		if md == nil {
			continue
		}
		if i >= size {
			s, q, err := storeMetadata(date, data)
			if err != nil {
				return 0, 0, err
			}
			affectedStock, affectedQuote = affectedStock+s, affectedQuote+q
			i = 0
			data = data[:0]
		} else {
			data = append(data, md)
		}
	}
	if len(data) > 0 {
		s, d, err := storeMetadata(date, data)
		if err != nil {
			return 0, 0, err
		}
		affectedStock, affectedQuote = affectedStock+s, affectedQuote+d
	}

	return affectedStock, affectedQuote, nil
}

func storeMetadata(date time.Time, metadata []*model.Metadata) (int64, int64, error) {
	var (
		stocks = make([]*db.Stock, 0, len(metadata))
		days   = make([]*db.Quote, 0, len(metadata))

		affectedStock, affectedQuote int64
	)

	for _, md := range metadata {
		stocks = append(stocks, &db.Stock{
			Code:            md.Code,
			Name:            md.Name,
			Suspend:         md.Suspend,
			CreateTimestamp: time.Now(),
		})

		day, err := BuildQuoteDayWitchMetadata(md, date)
		if err != nil {
			return 0, 0, err
		}
		days = append(days, day)
	}

	if len(stocks) != 0 {
		affected, err := storeStock(stocks)
		if err != nil {
			return 0, 0, err
		}
		affectedStock += affected

		affected, err = storeQuote(days, db.Day, date)
		if err != nil {
			return 0, 0, err
		}
		affectedQuote += affected
	}

	if date.Weekday() == time.Friday {
		var offset, limit int64 = 0, 30

		for {
			stocks, err := db.StockWithSelectRange(mysql.DB, offset, limit, DBExecTimeout)
			if err != nil {
				return 0, 0, err
			}

			weeks := make([]*db.Quote, 0, len(stocks))
			for _, stock := range stocks {
				week, err := BuildQuoteWeekWithQuoteDay(stock.Code, date)
				if err != nil {
					return affectedStock, affectedQuote, err
				}
				if week != nil {
					weeks = append(weeks, week)
				}
			}

			if _, err := storeQuote(weeks, db.Week, date); err != nil {
				return affectedStock, affectedQuote, err
			}

			if int64(len(stocks)) < limit {
				break
			}
			offset += limit
		}

	}
	return affectedStock, affectedQuote, nil
}

func storeStock(data []*db.Stock) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	tx, err := mysql.DB.Begin()
	if err != nil {
		return 0, err
	}
	affected, err := db.StockWithInsertOrUpdateMany(tx, data, DBExecTimeout)
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

func storeQuote(data []*db.Quote, kind string, date time.Time) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if kind != db.Day && kind != db.Week {
		return 0, fmt.Errorf("invalid kind: %v", kind)
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
	if _, err := db.QuoteWithDeleteManyByCodesAndDate(tx, kind, codes, d, DBExecTimeout); err != nil {
		tx.Rollback()
		return 0, err
	}

	affected, err := db.QuoteWithInsertMany(tx, kind, data, DBExecTimeout)
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
