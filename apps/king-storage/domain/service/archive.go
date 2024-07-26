package service

import (
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-storage/domain/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/model"
)

var timeout = 30 * time.Second

func ArchiveMetadata(date time.Time, metadata chan *model.Metadata) (int64, int64, int64, error) {
	var (
		affectedStock, affectedQuoteDay, affectedQuoteWeek int64
		i, size                                            = 0, 50
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
			affectedStock, affectedQuoteDay, affectedQuoteWeek = affectedStock+s, affectedQuoteDay+d, affectedQuoteWeek+w
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
		affectedStock, affectedQuoteDay, affectedQuoteWeek = affectedStock+s, affectedQuoteDay+d, affectedQuoteWeek+w
	}

	return affectedStock, affectedQuoteDay, affectedQuoteWeek, nil
}

func storeMetadata(date time.Time, metadata []*model.Metadata) (int64, int64, int64, error) {
	var (
		stocks = make([]*db.Stock, 0, len(metadata))
		days   = make([]*db.Quote, 0, len(metadata))

		affectedStocktock, affectedQuoteDay, affectedQuoteWeek int64
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
			return 0, 0, 0, err
		}
		days = append(days, day)
	}

	if len(stocks) != 0 {
		affected, err := storeStock(stocks)
		if err != nil {
			return 0, 0, 0, err
		}
		affectedStocktock += affected

		affected, err = storeQuote(days, db.Day, date)
		if err != nil {
			return 0, 0, 0, err
		}
		affectedQuoteDay += affected
	}

	if date.Weekday() == time.Friday {
		var (
			offset, limit int64 = 0, 50
			timeout             = 10 * time.Second
		)

		for {
			stocks, err := db.StockWithSelectRange(mysql.DB, offset, limit, timeout)
			if err != nil {
				return 0, 0, 0, err
			}

			weeks := make([]*db.Quote, 0, len(stocks))
			for _, stock := range stocks {
				week, err := BuildQuoteWeekWithQuoteDay(stock.Code, date)
				if err != nil && err != ErrNoData {
					return affectedStocktock, affectedQuoteDay, affectedQuoteWeek, err
				}
				if err == ErrNoData {
					continue
				}
				if week != nil {
					weeks = append(weeks, week)
				}
			}

			affected, err := storeQuote(weeks, db.Week, date)
			if err != nil {
				return affectedStocktock, affectedQuoteDay, affectedQuoteWeek, err
			}
			affectedQuoteWeek += affected

			if int64(len(stocks)) < limit {
				break
			}
			offset += limit
		}

	}
	return affectedStocktock, affectedQuoteDay, affectedQuoteWeek, nil
}

func storeStock(data []*db.Stock) (int64, error) {
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
	if _, err := db.QuoteWithDeleteManyByCodesAndDate(tx, kind, codes, d, timeout); err != nil {
		tx.Rollback()
		return 0, err
	}

	affected, err := db.QuoteWithInsertMany(tx, kind, data, timeout)
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
