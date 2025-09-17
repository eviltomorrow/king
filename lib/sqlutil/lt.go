package sqlutil

import (
	"fmt"
)

func WithLt(expr string, value interface{}) condition {
	return &lt{prefix: AND, expr: expr, value: value}
}

type lt struct {
	prefix Type
	expr   string
	value  interface{}
}

func (c *lt) SQL() (string, []interface{}) {
	return c.prefix.String() + fmt.Sprintf("%s%s < ?", c.prefix.String(), c.expr), []interface{}{c.value}
}

func (c *lt) SetPrefix(prefix Type) condition {
	c.prefix = prefix
	return c
}
