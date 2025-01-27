package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-storage/domain/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/model"
)

func StoreMetadata(ctx context.Context, metadata []*model.Metadata, date time.Time) (int64, int64, int64, error) {
	if len(metadata) == 0 {
		return 0, 0, 0, nil
	}

	var affectedStock, affectedDay, affectedWeek int64

	stocks, err := BuildStocksWithMetadata(ctx, metadata)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(stocks) != 0 {
		affected, err := storeStock(ctx, stocks)
		if err != nil {
			return 0, 0, 0, err
		}
		affectedStock += affected
	}

	days, err := BuildQuoteDaysWitchMetadata(ctx, metadata, date)
	if err != nil {
		return 0, 0, 0, err
	}
	if len(days) != 0 {
		affected, err := storeQuote(ctx, days, db.Day, date)
		if err != nil {
			return 0, 0, 0, err
		}
		affectedDay += affected
	}

	if date.Weekday() == time.Friday {
		weeks, err := BuildQuoteWeeksWithMetadata(ctx, metadata, date)
		if err != nil {
			return 0, 0, 0, err
		}
		if len(weeks) != 0 {
			affected, err := storeQuote(ctx, weeks, db.Week, date)
			if err != nil {
				return 0, 0, 0, err
			}
			affectedWeek += affected
		}
	}
	return affectedStock, affectedDay, affectedWeek, nil
}

func storeStock(ctx context.Context, data []*db.Stock) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	tx, err := mysql.DB.Begin()
	if err != nil {
		return 0, err
	}
	affected, err := db.StockWithInsertOrUpdateMany(ctx, tx, data)
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

func storeQuote(ctx context.Context, data []*db.Quote, kind string, date time.Time) (int64, error) {
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
	if _, err := db.QuoteWithDeleteManyByCodesAndDate(ctx, tx, kind, codes, d); err != nil {
		tx.Rollback()
		return 0, err
	}

	affected, err := db.QuoteWithInsertMany(ctx, tx, kind, data)
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
