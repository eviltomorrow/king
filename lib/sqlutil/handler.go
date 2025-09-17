package sqlutil

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

type handler struct {
	exec mysql.Exec

	builder strings.Builder
	args    []interface{}

	table      string
	columns    []string
	conditions []condition
	groupby    []string
	orderby    []string
	limit      []int64
}

func (h *handler) QueryOne(f func(row *sql.Row) error) error {
	return nil
}

func (h *handler) Query(f func(row *sql.Rows) error) error {
	column := "*"
	if len(h.columns) != 0 {
		column = strings.Join(h.columns, ", ")
	}

	table := h.table

	var (
		where strings.Builder
		args  []interface{}
	)
	for _, condition := range h.conditions {
		c, arg := condition.SQL()
		if _, err := where.WriteString(c); err != nil {
			return err
		}
		args = append(args, arg...)
	}

	groupby := ""
	if len(h.groupby) != 0 {
		groupby = strings.Join(h.groupby, ",")
	}

	orderby := ""
	if len(h.orderby) != 0 {
		orderby = strings.Join(h.orderby, ",")
	}

	sql := fmt.Sprintf("SELECT %s FROM %s", column, table)
	if where.Len() != 0 {
		w := where.String()

		w = strings.TrimPrefix(w, AND.String())
		w = strings.TrimPrefix(w, OR.String())
		sql = fmt.Sprintf("%s WHERE %s", sql, w)
	}
	if groupby != "" {
		sql = fmt.Sprintf("%s GROUP BY %s", sql, groupby)
	}
	if orderby != "" {
		sql = fmt.Sprintf("%s ORDER BY %s", sql, orderby)
	}
	if len(h.limit) == 2 {
		sql = fmt.Sprintf("%s LIMIT ?, ?", sql)
		args = append(args, h.limit[0], h.limit[1])
	}

	fmt.Println(sql, args)

	return nil
}

func (h *handler) Update(map[string]interface{}) (int64, error) {
	return 0, nil
}

func (h *handler) Delete() (int64, error) {
	return 0, nil
}

func (h *handler) Insert(map[string]interface{}) (int64, error) {
	return 0, nil
}

func (h *handler) Columns(columns []string) QueryHandler {
	h.columns = columns
	return h
}

func (h *handler) Table(name string) QueryHandler {
	h.table = name
	return h
}

func (h *handler) Where(conditions ...condition) QueryHandler {
	h.conditions = append(h.conditions, conditions...)
	return h
}

func (h *handler) GroupBy(gropuBy []string) QueryHandler {
	h.groupby = gropuBy
	return h
}

func (h *handler) OrderBy(orderBy ...order) QueryHandler {
	for _, o := range orderBy {
		text, args := o.SQL()
		h.orderby = append(h.orderby, text)
		h.args = append(h.args, args...)
	}
	return h
}

func (h *handler) Limit(offset int64, count int64) QueryHandler {
	h.limit = []int64{offset, count}
	return h
}
