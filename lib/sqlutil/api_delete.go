package sqlutil

import (
	"strings"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func NewDeleteHandler(exec mysql.Exec) DeleteHandler {
	return &handler{
		exec:    exec,
		builder: strings.Builder{},
		args:    make([]interface{}, 0, 16),
	}
}

type DeleteHandler interface {
	Delete() (int64, error)

	DeleteTable
	DeleteWhere
}

type DeleteTable interface {
	Table(name string) QueryHandler
}

type DeleteWhere interface {
	Where(...condition) QueryHandler
}
