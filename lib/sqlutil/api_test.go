package sqlutil

import (
	"testing"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func TestQuery(t *testing.T) {
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

func TestDelete(t *testing.T) {
	NewDelete(mysql.DB).Table("C").Where(WithBetweenAnd("c", 10, 20)).Delete()
}

func TestInsert(t *testing.T) {
	NewInsert(mysql.DB).Table("C").Insert(map[string]interface{}{"c": "is"})
}

func TestInsertBatch(t *testing.T) {
	NewInsert(mysql.DB).Table("C").InsertBatch([]string{"c"}, []map[string]interface{}{{"c": "is"}, {"c": "is"}})
}

func TestUpdate(t *testing.T) {
	NewUpdate(mysql.DB).Table("DB").Field(map[string]interface{}{"c": "d"}).Update()
}
