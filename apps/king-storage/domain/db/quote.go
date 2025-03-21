package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/mathutil"
	"github.com/eviltomorrow/king/lib/sqlutil"
	jsoniter "github.com/json-iterator/go"
)

const (
	Day  = "day"
	Week = "week"
)

func QuoteWithInsertMany(ctx context.Context, exec mysql.Exec, kind string, data []*Quote) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	FieldQuotes := make([]string, 0, len(data))
	args := make([]interface{}, 0, 12*len(data))
	for _, m := range data {
		FieldQuotes = append(FieldQuotes, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now())")
		args = append(args, m.Id)
		args = append(args, m.Code)
		args = append(args, m.Open)
		args = append(args, m.Close)
		args = append(args, m.High)
		args = append(args, m.Low)
		args = append(args, m.YesterdayClosed)
		args = append(args, m.Volume)
		args = append(args, m.Account)
		args = append(args, m.Date.Format(time.DateOnly))
		args = append(args, m.NumOfYear)
		args = append(args, m.Xd)
	}

	_sql := fmt.Sprintf("insert into quote_%s (%s) values %s", kind, strings.Join(quoteFeilds, ","), strings.Join(FieldQuotes, ","))
	result, err := exec.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func QuoteWithDeleteManyByCodesAndDate(ctx context.Context, exec mysql.Exec, kind string, codes []string, date string) (int64, error) {
	if len(codes) == 0 {
		return 0, nil
	}

	FieldQuotes := make([]string, 0, len(codes))
	args := make([]interface{}, 0, len(codes)+1)
	for _, code := range codes {
		FieldQuotes = append(FieldQuotes, "?")
		args = append(args, code)
	}
	args = append(args, date)

	_sql := fmt.Sprintf("delete from quote_%s where code in (%s) and DATE_FORMAT(`date`, '%%Y-%%m-%%d') = ?", kind, strings.Join(FieldQuotes, ","))
	result, err := exec.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func QuoteWithSelectBetweenByCodeAndDate(ctx context.Context, exec mysql.Exec, kind string, code string, begin, end string) ([]*Quote, error) {
	_sql := fmt.Sprintf("select id, code, open, close, high, low, yesterday_closed, volume, account, date, num_of_year, xd, create_timestamp, modify_timestamp from quote_%s where code = ? and DATE_FORMAT(`date`, '%%Y-%%m-%%d') between ? and ? order by date asc", kind)

	rows, err := exec.QueryContext(ctx, _sql, code, begin, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]*Quote, 0, 5)
	for rows.Next() {
		m := Quote{}
		if err := rows.Scan(
			&m.Id,
			&m.Code,
			&m.Open,
			&m.Close,
			&m.High,
			&m.Low,
			&m.YesterdayClosed,
			&m.Volume,
			&m.Account,
			&m.Date,
			&m.NumOfYear,
			&m.Xd,
			&m.CreateTimestamp,
			&m.ModifyTimestamp,
		); err != nil {
			return nil, err
		}
		data = append(data, &m)
	}

	result := make([]*Quote, len(data))
	var xd float64 = 1.0
	for i := len(data) - 1; i >= 0; i-- {
		d := data[i]
		if xd != 1.0 {
			n := &Quote{
				Id:              d.Id,
				Code:            d.Code,
				Open:            mathutil.Trunc2(d.Open * xd),
				Close:           mathutil.Trunc2(d.Close * xd),
				High:            mathutil.Trunc2(d.High * xd),
				Low:             mathutil.Trunc2(d.Low * xd),
				YesterdayClosed: mathutil.Trunc2(d.YesterdayClosed * xd),
				Volume:          d.Volume,
				Account:         d.Account,
				Date:            d.Date,
				NumOfYear:       d.NumOfYear,
				Xd:              d.Xd,
				CreateTimestamp: d.CreateTimestamp,
				ModifyTimestamp: d.ModifyTimestamp,
			}
			result[i] = n
		} else {
			result[i] = d
		}

		if d.Xd != 1.0 {
			xd = xd * d.Xd
		}
	}

	return result, nil
}

func QuoteWithSelectBetweenByCodesAndDate(ctx context.Context, exec mysql.Exec, kind string, codes []string, begin, end string) (map[string][]*Quote, error) {
	if len(codes) == 0 {
		return nil, nil
	}
	fields := make([]string, 0, len(codes))
	args := make([]interface{}, 0, len(codes)+2)
	for _, code := range codes {
		fields = append(fields, "?")
		args = append(args, code)
	}
	args = append(args, begin)
	args = append(args, end)

	_sql := fmt.Sprintf("select id, code, open, close, high, low, yesterday_closed, volume, account, date, num_of_year, xd, create_timestamp, modify_timestamp from quote_%s where code in (%s) and DATE_FORMAT(`date`, '%%Y-%%m-%%d') between ? and ? order by date asc", kind, strings.Join(fields, ","))
	rows, err := exec.QueryContext(ctx, _sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[string][]*Quote, len(codes))
	for rows.Next() {
		m := Quote{}
		if err := rows.Scan(
			&m.Id,
			&m.Code,
			&m.Open,
			&m.Close,
			&m.High,
			&m.Low,
			&m.YesterdayClosed,
			&m.Volume,
			&m.Account,
			&m.Date,
			&m.NumOfYear,
			&m.Xd,
			&m.CreateTimestamp,
			&m.ModifyTimestamp,
		); err != nil {
			return nil, err
		}
		v, ok := data[m.Code]
		if !ok {
			v = make([]*Quote, 0, 5)
		}
		v = append(v, &m)
		data[m.Code] = v
	}

	for _, v := range data {
		var xd float64 = 1.0
		for i := len(v) - 1; i >= 0; i-- {
			d := v[i]
			if xd != 1.0 {
				n := &Quote{
					Id:              d.Id,
					Code:            d.Code,
					Open:            mathutil.Trunc2(d.Open * xd),
					Close:           mathutil.Trunc2(d.Close * xd),
					High:            mathutil.Trunc2(d.High * xd),
					Low:             mathutil.Trunc2(d.Low * xd),
					YesterdayClosed: mathutil.Trunc2(d.YesterdayClosed * xd),
					Volume:          d.Volume,
					Account:         d.Account,
					Date:            d.Date,
					NumOfYear:       d.NumOfYear,
					Xd:              d.Xd,
					CreateTimestamp: d.CreateTimestamp,
					ModifyTimestamp: d.ModifyTimestamp,
				}
				v[i] = n
			} else {
				v[i] = d
			}

			if d.Xd != 1.0 {
				xd = xd * d.Xd
			}
		}
	}

	return data, nil
}

func QuoteWithSelectLatestByCodesAndDate(ctx context.Context, exec mysql.Exec, kind string, code []string, date string) ([]*Quote, error) {
	if len(code) == 0 {
		return nil, nil
	}
	fields := make([]string, 0, len(code))
	args := make([]interface{}, 0, len(code)+1)
	for _, c := range code {
		fields = append(fields, "?")
		args = append(args, c)
	}
	args = append(args, date)

	_sql := fmt.Sprintf("select d1.id, d1.code, d1.open, d1.close, d1.high, d1.low, d1.yesterday_closed, d1.volume, d1.account, d1.date, d1.num_of_year, d1.xd, d1.create_timestamp, d1.modify_timestamp from quote_%s d1 inner join (select code, max(date) `date` from quote_day where code in (%s) and DATE_FORMAT(`date`, '%%Y-%%m-%%d') <= ? group by code) d2 on d1.code = d2.code and d1.`date` = d2.`date`", kind, strings.Join(fields, ","))
	data := make([]*Quote, 0, len(code))
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			m := Quote{}
			if err := rows.Scan(
				&m.Id,
				&m.Code,
				&m.Open,
				&m.Close,
				&m.High,
				&m.Low,
				&m.YesterdayClosed,
				&m.Volume,
				&m.Account,
				&m.Date,
				&m.NumOfYear,
				&m.Xd,
				&m.CreateTimestamp,
				&m.ModifyTimestamp,
			); err != nil {
				return err
			}
			data = append(data, &m)
		}
		return nil
	}

	if err := sqlutil.ExecQueryMany(ctx, exec, _sql, args, scan); err != nil {
		return nil, err
	}

	return data, nil
}

func QuoteWithSelectLatestByCodeAndDate(ctx context.Context, exec mysql.Exec, kind string, code string, date string, limit int64) ([]*Quote, error) {
	_sql := fmt.Sprintf("select id, code, open, close, high, low, yesterday_closed, volume, account, date, num_of_year, xd, create_timestamp, modify_timestamp from quote_%s where code = ? and DATE_FORMAT(`date`, '%%Y-%%m-%%d') <= ? order by `date` desc limit ?", kind)
	rows, err := exec.QueryContext(ctx, _sql, code, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]*Quote, 0, limit)
	for rows.Next() {
		m := Quote{}
		if err := rows.Scan(
			&m.Id,
			&m.Code,
			&m.Open,
			&m.Close,
			&m.High,
			&m.Low,
			&m.YesterdayClosed,
			&m.Volume,
			&m.Account,
			&m.Date,
			&m.NumOfYear,
			&m.Xd,
			&m.CreateTimestamp,
			&m.ModifyTimestamp,
		); err != nil {
			return nil, err
		}
		data = append(data, &m)
	}

	result := make([]*Quote, 0, len(data))
	xd := 1.0
	for _, d := range data {
		if xd != 1.0 {
			n := &Quote{
				Id:              d.Id,
				Code:            d.Code,
				Open:            mathutil.Trunc2(d.Open * xd),
				Close:           mathutil.Trunc2(d.Close * xd),
				High:            mathutil.Trunc2(d.High * xd),
				Low:             mathutil.Trunc2(d.Low * xd),
				YesterdayClosed: mathutil.Trunc2(d.YesterdayClosed * xd),
				Volume:          d.Volume,
				Account:         d.Account,
				Date:            d.Date,
				NumOfYear:       d.NumOfYear,
				Xd:              d.Xd,
				CreateTimestamp: d.CreateTimestamp,
				ModifyTimestamp: d.ModifyTimestamp,
			}
			result = append(result, n)
		} else {
			result = append(result, d)
		}

		if d.Xd != 1.0 {
			xd = d.Xd
		}
	}

	return result, nil
}

func QuoteWithSelectRangeByDate(ctx context.Context, exec mysql.Exec, kind string, date string, offset, limit int64) ([]*Quote, error) {
	_sql := fmt.Sprintf("select id, code, open, close, high, low, yesterday_closed, volume, account, date, num_of_year, xd, create_timestamp, modify_timestamp from quote_%s where DATE_FORMAT(`date`, '%%Y-%%m-%%d') = ? limit ?, ?", kind)
	rows, err := exec.QueryContext(ctx, _sql, date, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]*Quote, 0, limit)
	for rows.Next() {
		m := Quote{}
		if err := rows.Scan(
			&m.Id,
			&m.Code,
			&m.Open,
			&m.Close,
			&m.High,
			&m.Low,
			&m.YesterdayClosed,
			&m.Volume,
			&m.Account,
			&m.Date,
			&m.NumOfYear,
			&m.Xd,
			&m.CreateTimestamp,
			&m.ModifyTimestamp,
		); err != nil {
			return nil, err
		}
		data = append(data, &m)
	}

	return data, nil
}

func QuoteWithSelectOneByCodeAndDate(ctx context.Context, exec mysql.Exec, kind string, code string, date string) (*Quote, error) {
	_sql := fmt.Sprintf("select id, code, open, close, high, low, yesterday_closed, volume, account, date, num_of_year, xd, create_timestamp, modify_timestamp from quote_%s where code = ? and DATE_FORMAT(`date`, '%%Y-%%m-%%d') = ?", kind)
	row := exec.QueryRowContext(ctx, _sql, code, date)
	if row.Err() != nil {
		return nil, row.Err()
	}
	m := Quote{}
	if err := row.Scan(
		&m.Id,
		&m.Code,
		&m.Open,
		&m.Close,
		&m.High,
		&m.Low,
		&m.YesterdayClosed,
		&m.Volume,
		&m.Account,
		&m.Date,
		&m.NumOfYear,
		&m.Xd,
		&m.CreateTimestamp,
		&m.ModifyTimestamp,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func QuoteWithCountByDate(ctx context.Context, exec mysql.Exec, kind string, date string) (int64, error) {
	return sqlutil.TableWithCount(ctx, exec, fmt.Sprintf("quote_%s", kind), map[string]interface{}{FieldQuoteDate: date})
}

type Quote struct {
	Id              string       `json:"id"`
	Code            string       `json:"code"`
	Open            float64      `json:"open"`
	Close           float64      `json:"close"`
	High            float64      `json:"high"`
	Low             float64      `json:"low"`
	YesterdayClosed float64      `json:"yesterday_closed"`
	Volume          int64        `json:"volume"`
	Account         float64      `json:"account"`
	Date            time.Time    `json:"date"`
	NumOfYear       int          `json:"num_of_year"`
	Xd              float64      `json:"xd"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (q *Quote) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}

const (
	FieldQuoteID              = "id"
	FieldQuoteCode            = "code"
	FieldQuoteOpen            = "open"
	FieldQuoteClose           = "close"
	FieldQuoteHigh            = "high"
	FieldQuoteLow             = "low"
	FieldQuoteYesterdayClosed = "yesterday_closed"
	FieldQuoteVolume          = "volume"
	FieldQuoteAccount         = "account"
	FieldQuoteDate            = "date"
	FieldQuoteNumOfYear       = "num_of_year"
	FieldQuoteXd              = "xd"
	FieldQuoteCreateTimestamp = "create_timestamp"
	FieldQuoteModifyTimestamp = "modify_timestamp"
)

var quoteFeilds = []string{
	FieldQuoteID,
	FieldQuoteCode,
	FieldQuoteOpen,
	FieldQuoteClose,
	FieldQuoteHigh,
	FieldQuoteLow,
	FieldQuoteYesterdayClosed,
	FieldQuoteVolume,
	FieldQuoteAccount,
	FieldQuoteDate,
	FieldQuoteNumOfYear,
	FieldQuoteXd,
	FieldQuoteCreateTimestamp,
}
