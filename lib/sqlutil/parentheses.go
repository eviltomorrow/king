package sqlutil

import (
	"fmt"
	"strings"
)

func WithParentheses(conditions ...condition) condition {
	return &parentheses{prefix: AND, conditions: conditions}
}

type parentheses struct {
	prefix     Type
	conditions []condition
}

func (c *parentheses) SQL() (string, []interface{}) {
	var (
		sql    strings.Builder
		values = make([]interface{}, 0, 32)
	)
	for _, condition := range c.conditions {
		text, args := condition.SQL()
		sql.WriteString(text)
		values = append(values, args...)
	}
	w := sql.String()
	w = strings.TrimPrefix(w, AND.String())
	w = strings.TrimPrefix(w, OR.String())
	return fmt.Sprintf("%s(%s)", c.prefix.String(), w), values
}

func (c *parentheses) SetPrefix(prefix Type) condition {
	c.prefix = prefix
	return c
}
