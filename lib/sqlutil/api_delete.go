package sqlutil

import (
	"fmt"
	"strings"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func NewDelete(exec mysql.Exec) Delete {
	return &delete{
		exec:    exec,
		builder: strings.Builder{},
	}
}

type Delete interface {
	Delete() (int64, error)

	DeleteTable
	DeleteWhere
}

type DeleteTable interface {
	Table(name string) Delete
}

type DeleteWhere interface {
	Where(...condition) Delete
}

type delete struct {
	exec mysql.Exec

	builder strings.Builder
	args    []interface{}

	table      string
	conditions []condition
}

func (h *delete) Delete() (int64, error) {
	if h.table == "" {
		return 0, fmt.Errorf("table is invalid")
	}
	table := h.table

	var (
		where strings.Builder
		args  []interface{}
	)
	for _, condition := range h.conditions {
		c, arg := condition.SQL()
		if _, err := where.WriteString(c); err != nil {
			return 0, err
		}
		args = append(args, arg...)
	}

	sql := fmt.Sprintf("DELETE FROM %s", table)
	if where.Len() != 0 {
		w := where.String()

		w = strings.TrimPrefix(w, AND.String())
		w = strings.TrimPrefix(w, OR.String())
		sql = fmt.Sprintf("%s WHERE %s", sql, w)
	}

	fmt.Println(sql, args)

	return 0, nil
}

func (h *delete) Table(name string) Delete {
	h.table = name
	return h
}

func (h *delete) Where(conditions ...condition) Delete {
	h.conditions = conditions
	return h
}
