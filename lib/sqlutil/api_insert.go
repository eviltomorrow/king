package sqlutil

import (
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
	Insert() (int64, error)

	InsertField
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

	table  string
	fields map[string]interface{}
}

func (h *insert) Insert() (int64, error) {
	return 0, nil
}

func (h *insert) Field(fields map[string]interface{}) Insert {
	h.fields = fields
	return h
}

func (h *insert) Table(name string) Insert {
	return h
}
