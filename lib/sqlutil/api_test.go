package sqlutil

import (
	"testing"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func TestQuery(t *testing.T) {
	// Table("").Columns([]string{}).Where(WithBetweenAnd("", "", "")).GroupBy([]string{}).OrderBy().Query(nil)

	NewQuery(mysql.DB).
		Columns([]string{
			"A as c",
			"B",
			"D as D",
			"Date(C)",
		}).
		Table("C").
		Where(
			WithEq("C", 10),
			WithGt("D", "C").SetPrefix(OR),
			WithParentheses(
				WithEq("A", 10),
				WithGt("B", 20),
			).SetPrefix(OR),
			WithNotNull("C"),
		).
		GroupBy([]string{
			"ID",
		}).
		OrderBy(
			ASC("C"),
			ASC("D"),
			DESC("BB"),
		).
		Limit(10, 10).
		Query(nil)
}
