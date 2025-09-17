package sqlutil

import (
	"strings"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func NewInsertHandler(exec mysql.Exec) InsertHandler {
	return &handler{
		exec:    exec,
		builder: strings.Builder{},
		args:    make([]interface{}, 0, 16),
	}
}

type InsertHandler interface {
	Insert(map[string]interface{}) (int64, error)

	InsertTable
}

type InsertTable interface {
	Table(name string) QueryHandler
}
