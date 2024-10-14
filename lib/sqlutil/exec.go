package sqlutil

import (
	"context"
	"database/sql"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func ExecQueryMany(ctx context.Context, exec mysql.Exec, sql string, args []interface{}, f func(row *sql.Rows) error) error {
	rows, err := exec.QueryContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := f(rows); err != nil {
		return err
	}
	return rows.Err()
}

func ExecQueryOne(ctx context.Context, exec mysql.Exec, sql string, args []interface{}, f func(row *sql.Row) error) error {
	row := exec.QueryRowContext(ctx, sql, args...)
	if row.Err() != nil {
		return row.Err()
	}

	if err := f(row); err != nil {
		return err
	}
	return row.Err()
}
