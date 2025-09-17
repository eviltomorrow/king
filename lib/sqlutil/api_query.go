package sqlutil

import (
	"database/sql"
	"strings"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func NewQueryHandler(exec mysql.Exec) QueryHandler {
	return &handler{
		exec:    exec,
		builder: strings.Builder{},
		args:    make([]interface{}, 0, 16),
	}
}

type QueryHandler interface {
	QueryOne(f func(row *sql.Row) error) error
	Query(f func(row *sql.Rows) error) error

	QueryColumn
	QueryTable
	QueryWhere
	QueryGroupBy
	QueryOrderBy
	QueryLimit
}

type QueryColumn interface {
	Columns(columns []string) QueryHandler
}

type QueryTable interface {
	Table(name string) QueryHandler
}

type QueryWhere interface {
	Where(...condition) QueryHandler
}

type QueryGroupBy interface {
	GroupBy([]string) QueryHandler
}

type QueryOrderBy interface {
	OrderBy(...order) QueryHandler
}

type QueryLimit interface {
	Limit(int64, int64) QueryHandler
}
