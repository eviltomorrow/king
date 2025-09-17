package sqlutil

import "fmt"

func WithNotEq(expr string, value interface{}) condition {
	return &notEq{prefix: AND, expr: expr, value: value}
}

type notEq struct {
	prefix Type
	expr   string
	value  interface{}
}

func (c *notEq) SQL() (string, []interface{}) {
	return fmt.Sprintf("%s%s != ?", c.prefix.String(), c.expr), []interface{}{c.value}
}

func (c *notEq) SetPrefix(prefix Type) condition {
	c.prefix = prefix
	return c
}
