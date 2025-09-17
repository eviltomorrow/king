package sqlutil

import (
	"fmt"
)

func WithEq(expr string, value interface{}) condition {
	return &eq{prefix: AND, expr: expr, value: value}
}

type eq struct {
	prefix Type
	expr   string
	value  interface{}
}

func (c *eq) SQL() (string, []interface{}) {
	return fmt.Sprintf("%s%s = ?", c.prefix.String(), c.expr), []interface{}{c.value}
}

func (c *eq) SetPrefix(prefix Type) condition {
	c.prefix = prefix
	return c
}
