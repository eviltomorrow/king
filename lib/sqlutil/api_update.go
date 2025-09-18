package sqlutil

import (
	"strings"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func NewUpdate(exec mysql.Exec) Update {
	return &update{
		exec:    exec,
		builder: strings.Builder{},
	}
}

type Update interface {
	Update() (int64, error)

	UpdateField
	UpdateTable
	UpdateWhere
}

type UpdateField interface {
	Field(map[string]interface{}) Update
}

type UpdateTable interface {
	Table(name string) Update
}

type UpdateWhere interface {
	Where(...condition) Update
}

type update struct {
	exec mysql.Exec

	builder strings.Builder

	fields     map[string]interface{}
	table      string
	conditions []condition
}

func (h *update) Field(fields map[string]interface{}) Update {
	h.fields = fields
	return h
}

func (h *update) Table(name string) Update {
	h.table = name
	return h
}

func (h *update) Where(conditions ...condition) Update {
	h.conditions = append(h.conditions, conditions...)
	return h
}

func (h *update) Update() (int64, error) {
	return 0, nil
}
