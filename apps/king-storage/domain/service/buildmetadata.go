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

func BuildQuoteDayWitchMetadata(ctx context.Context, data *model.Metadata, date time.Time) (*db.Quote, error) {
	latest, err := db.QuoteWithSelectLatestByCodeAndDate(ctx, mysql.DB, db.Day, data.Code, data.Date, 1)
	if err != nil {
		return nil, err
	}

	var xd float64 = 1.0
	if len(latest) == 1 && latest[0].Close != 0 && latest[0].Date.Format(time.DateOnly) != data.Date && latest[0].Close != data.YesterdayClosed {
		xd = data.YesterdayClosed / latest[0].Close
	}

	quote := &db.Quote{
		Id:              snowflake.GenerateID(),
		Code:            data.Code,
		Open:            data.Open,
		Close:           data.Latest,
		High:            data.High,
		Low:             data.Low,
		YesterdayClosed: data.YesterdayClosed,
		Volume:          data.Volume,
		Account:         data.Account,
		Date:            date,
		NumOfYear:       date.YearDay(),
		Xd:              xd,
		CreateTimestamp: time.Now(),
	}
	return quote, nil
}

func BuildQuoteDaysWitchMetadata(ctx context.Context, data []*model.Metadata, date time.Time) ([]*db.Quote, error) {
	codes := make([]string, 0, len(data))
	for _, d := range data {
		codes = append(codes, d.Code)
	}

	latest, err := db.QuoteWithSelectLatestByCodesAndDate(ctx, mysql.DB, db.Day, codes, date.Format(time.DateOnly))
	if err != nil {
		return nil, err
	}

	md := make(map[string]*db.Quote, len(latest))
	for _, l := range latest {
		md[l.Code] = l
	}

	result := make([]*db.Quote, 0, len(data))
	for _, d := range data {
		var xd float64 = 1.0

		m, ok := md[d.Code]
		if ok {
			if m.Close != 0 && m.Date.Format(time.DateOnly) != d.Date && m.Close != d.YesterdayClosed {
				xd = d.YesterdayClosed / m.Close
			}
		}

		quote := &db.Quote{
			Id:              snowflake.GenerateID(),
			Code:            d.Code,
			Open:            d.Open,
			Close:           d.Latest,
			High:            d.High,
			Low:             d.Low,
			YesterdayClosed: d.YesterdayClosed,
			Volume:          d.Volume,
			Account:         d.Account,
			Date:            date,
			NumOfYear:       date.YearDay(),
			Xd:              xd,
			CreateTimestamp: time.Now(),
		}
		result = append(result, quote)
	}
	return result, nil
}

func BuildQuoteWeekWithQuoteDay(ctx context.Context, code string, date time.Time) (*db.Quote, error) {
	if date.Weekday() != time.Friday {
		return nil, fmt.Errorf("panic: date is not friday")
	}
	var (
		begin = date.AddDate(0, 0, -5).Format(time.DateOnly)
		end   = date.Format(time.DateOnly)
	)

	days, err := db.QuoteWithSelectBetweenByCodeAndDate(ctx, mysql.DB, db.Day, code, begin, end)
	if err != nil {
		return nil, err
	}

	if len(days) == 0 {
		return nil, nil
	}

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
	return week, nil
}
