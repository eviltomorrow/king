package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-storage/domain/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/mathutil"
	"github.com/eviltomorrow/king/lib/model"
	"github.com/eviltomorrow/king/lib/snowflake"
	"github.com/eviltomorrow/king/lib/timeutil"
)

func BuildQuoteDaysWithMetadata(ctx context.Context, data []*model.Metadata, date time.Time) ([]*db.Quote, error) {
	return nil, nil
}

func BuildQuoteWeeksWithMetadata(ctx context.Context, md []*model.Metadata, date time.Time) ([]*db.Quote, error) {
	if date.Weekday() != time.Friday {
		return nil, fmt.Errorf("panic: date is not friday")
	}
	var (
		begin = date.AddDate(0, 0, -5).Format(time.DateOnly)
		end   = date.Format(time.DateOnly)
	)

	var codes = make([]string, 0, len(md))
	for _, d := range md {
		codes = append(codes, d.Code)
	}
	data, err := db.QuoteWithSelectBetweenByCodesAndDate(ctx, mysql.DB, db.Day, codes, begin, end)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	var result = make([]*db.Quote, 0, len(data)*5)
	for _, days := range data {
		var (
			first, last = days[0], days[len(days)-1]
			highs       = make([]float64, 0, len(days))
			lows        = make([]float64, 0, len(days))
			volumes     = make([]uint64, 0, len(days))
			accounts    = make([]float64, 0, len(days))
		)

		xd := 1.0
		for _, d := range days {
			highs = append(highs, d.High)
			lows = append(lows, d.Low)
			volumes = append(volumes, d.Volume)
			accounts = append(accounts, d.Account)
			if d.Xd != 1.0 {
				xd = d.Xd
			}
		}

		week := &db.Quote{
			Id:              snowflake.GenerateID(),
			Code:            first.Code,
			Open:            first.Open,
			Close:           last.Close,
			High:            mathutil.Max(highs),
			Low:             mathutil.Min(lows),
			YesterdayClosed: first.YesterdayClosed,
			Volume:          mathutil.Sum(volumes),
			Account:         mathutil.Sum(accounts),
			Date:            date,
			NumOfYear:       timeutil.YearWeek(date),
			Xd:              xd,
			CreateTimestamp: time.Now(),
		}
		result = append(result, week)
	}

	return result, nil
}

func BuildStocksWithMetadata(ctx context.Context, md []*model.Metadata) ([]*db.Stock, error) {
	stocks := make([]*db.Stock, 0, len(md))
	for _, md := range md {
		stocks = append(stocks, &db.Stock{
			Code:            md.Code,
			Name:            md.Name,
			Suspend:         md.Suspend,
			CreateTimestamp: time.Now(),
		})
	}
	return stocks, nil
}
