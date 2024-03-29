package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	jsoniter "github.com/json-iterator/go"
)

func StockWithInsertOrUpdateMany(exec mysql.Exec, stocks []*Stock, timeout time.Duration) (int64, error) {
	if len(stocks) == 0 {
		return 0, nil
	}

	codes := make([]string, 0, len(stocks))
	for _, stock := range stocks {
		codes = append(codes, stock.Code)
	}

	data, err := StockWithSelectMany(exec, codes, timeout)
	if err != nil {
		return 0, err
	}

	shouldInsertStocks := make([]*Stock, 0, len(stocks))
	shouldUpdateStocks := make([]*Stock, 0, len(stocks))
	for _, stock := range stocks {
		d, ok := data[stock.Code]
		if !ok {
			shouldInsertStocks = append(shouldInsertStocks, stock)
		} else {
			if d.Name != stock.Name {
				shouldUpdateStocks = append(shouldUpdateStocks, stock)
			}
		}
	}

	var count int64
	for _, s := range shouldUpdateStocks {
		affected, err := StockWithUpdateOne(exec, s.Code, s, timeout)
		if err != nil {
			return 0, err
		}
		count += affected
	}

	affected, err := StockWithInsertMany(exec, shouldInsertStocks, timeout)
	if err != nil {
		return 0, err
	}
	count += affected

	return count, nil
}

func StockWithInsertMany(exec mysql.Exec, stocks []*Stock, timeout time.Duration) (int64, error) {
	if len(stocks) == 0 {
		return 0, nil
	}

	var (
		exist = make(map[string]struct{}, len(stocks))
		data  = make([]*Stock, 0, len(stocks))
	)
	for _, stock := range stocks {
		if _, ok := exist[stock.Code]; !ok {
			data = append(data, stock)
			exist[stock.Code] = struct{}{}
		}
	}

	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	fields := make([]string, 0, len(data))
	args := make([]interface{}, 0, 3*len(data))
	for _, record := range data {
		fields = append(fields, "(?, ?, ?, now(), null)")
		args = append(args, record.Code)
		args = append(args, record.Name)
		args = append(args, record.Suspend)
	}

	_sql := fmt.Sprintf("insert into stock (%s) values %s", strings.Join(stockFields, ","), strings.Join(fields, ","))
	result, err := exec.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func StockWithUpdateOne(exec mysql.Exec, code string, stock *Stock, timeout time.Duration) (int64, error) {
	if stock == nil {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	_sql := `update stock set name = ?, suspend = ?, modify_timestamp = now() where code = ?`
	result, err := exec.ExecContext(ctx, _sql, stock.Name, stock.Suspend, code)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func StockWithSelectMany(exec mysql.Exec, codes []string, timeout time.Duration) (map[string]*Stock, error) {
	if len(codes) == 0 {
		return map[string]*Stock{}, nil
	}
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	fields := make([]string, 0, len(codes))
	args := make([]interface{}, 0, len(codes))
	for _, code := range codes {
		fields = append(fields, "?")
		args = append(args, code)
	}

	_sql := fmt.Sprintf(`select code, name, suspend, create_timestamp, modify_timestamp from stock where code in (%s)`, strings.Join(fields, ","))
	rows, err := exec.QueryContext(ctx, _sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stocks := make(map[string]*Stock, len(codes))
	for rows.Next() {
		stock := &Stock{}
		if err := rows.Scan(&stock.Code, &stock.Name, &stock.Suspend, &stock.CreateTimestamp, &stock.ModifyTimestamp); err != nil {
			return nil, err
		}
		stocks[stock.Code] = stock
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return stocks, nil
}

func StockWithSelectRange(exec mysql.Exec, offset, limit int64, timeout time.Duration) ([]*Stock, error) {
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	_sql := `select code, name, suspend, create_timestamp, modify_timestamp from stock limit ?, ?`
	rows, err := exec.QueryContext(ctx, _sql, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stocks := make([]*Stock, 0, limit)
	for rows.Next() {
		stock := &Stock{}
		if err := rows.Scan(&stock.Code, &stock.Name, &stock.Suspend, &stock.CreateTimestamp, &stock.ModifyTimestamp); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return stocks, nil
}

// Stock
type Stock struct {
	Code            string       `json:"code"`
	Name            string       `json:"name"`
	Suspend         string       `json:"suspend"`
	CreateTimestamp time.Time    `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime `json:"modify_timestamp"`
}

func (s *Stock) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(s)
	return string(buf)
}

const (
	FieldStockCode            = "code"
	FieldStockName            = "name"
	FieldStockSuspend         = "suspend"
	FieldStockCreateTimestamp = "create_timestamp"
	FieldStockModifyTimestamp = "modify_timestamp"
)

var stockFields = []string{
	FieldStockCode,
	FieldStockName,
	FieldStockSuspend,
	FieldStockCreateTimestamp,
	FieldStockModifyTimestamp,
}
