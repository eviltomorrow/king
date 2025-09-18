package sqlutil

import (
	"fmt"
	"strings"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func NewInsert(exec mysql.Exec) Insert {
	return &insert{
		exec:    exec,
		builder: strings.Builder{},
		args:    make([]interface{}, 0, 16),
	}
}

type Insert interface {
	Insert(map[string]interface{}) (int64, error)
	InsertBatch([]string, []map[string]interface{}) (int64, error)

	InsertTable
}

type InsertField interface {
	Field(map[string]interface{}) Insert
}

type InsertTable interface {
	Table(name string) Insert
}

type insert struct {
	exec mysql.Exec

	builder strings.Builder
	args    []interface{}

	table string
}

func (h *insert) Insert(value map[string]interface{}) (int64, error) {
	if h.table == "" {
		return 0, fmt.Errorf("table is invalid")
	}
	if len(value) == 0 {
		return 0, fmt.Errorf("value is invalid")
	}

	table := h.table
	fields1 := make([]string, 0, len(value))
	fields2 := make([]string, 0, len(value))
	args := make([]interface{}, 0, len(value))

	for k, v := range value {
		fields1 = append(fields1, k)
		fields2 = append(fields2, "?")
		args = append(args, v)
	}

	sql := fmt.Sprintf(`INSERT INTO %s(%s) VALUES (%s)`, table, strings.Join(fields1, ", "), strings.Join(fields2, ", "))

	fmt.Println(sql, args)
	return 0, nil
}

func (h *insert) InsertBatch(columns []string, values []map[string]interface{}) (int64, error) {
	if h.table == "" {
		return 0, fmt.Errorf("table is invalid")
	}
	if len(columns) == 0 {
		return 0, fmt.Errorf("columns is invalid")
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("values is invalid")
	}
	table := h.table

	fields := make([]string, 0, len(columns))
	for range columns {
		fields = append(fields, "?")
	}
	field := fmt.Sprintf("(%s)", strings.Join(fields, ", "))

	args := make([]interface{}, 0, len(columns)*len(values))
	fields = make([]string, 0, len(values))
	for _, value := range values {
		for _, column := range columns {
			arg, ok := value[column]
			if !ok {
				return 0, fmt.Errorf("not found column with [%s]", column)
			}
			args = append(args, arg)
		}
		fields = append(fields, field)
	}

	sql := fmt.Sprintf("insert into %s (%s) values %s", table, strings.Join(columns, ","), strings.Join(fields, ","))
	fmt.Println(sql, args)
	return 0, nil
}

func (h *insert) Table(name string) Insert {
	h.table = name
	return h
}
