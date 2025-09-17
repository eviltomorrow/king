package sqlutil

import (
	"strings"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func NewUpdateHandler(exec mysql.Exec) QueryHandler {
	return &handler{
		exec:    exec,
		builder: strings.Builder{},
		args:    make([]interface{}, 0, 16),
	}
}

type UpdateHandler interface {
	Update(map[string]interface{}) (int64, error)

	UpdateTable
	UpdateWhere
}

type UpdateTable interface {
	Table(name string) QueryHandler
}

type UpdateWhere interface {
	Where(...condition) QueryHandler
}
