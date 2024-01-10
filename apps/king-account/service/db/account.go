package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func AccountWithSelectOne(exec mysql.Exec, id string, timeout time.Duration) (*Account, error) {
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	_sql := `select id, username, password, nick_name, phone, email, status, register_datetime, create_timestamp, modify_timestamp from account where id = ?`
	row := exec.QueryRowContext(ctx, _sql, id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	a := Account{}
	if err := row.Scan(
		&a.Id,
		&a.Username,
		&a.Password,
		&a.NickName,
		&a.Phone,
		&a.Email,
		&a.Status,
		&a.RegisterDatetime,
		&a.CreateTimestamp,
		&a.ModifyTimestamp,
	); err != nil {
		return nil, err
	}
	return &a, nil
}

func AccountWithSelectRange(exec mysql.Exec, offset, limit int64, timeout time.Duration) ([]*Account, error) {
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	_sql := `select id, username, password, nick_name, phone, email, status, register_datetime, create_timestamp, modify_timestamp from account limit ?, ?`
	rows, err := exec.QueryContext(ctx, _sql, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]*Account, 0, limit)
	for rows.Next() {
		a := Account{}
		if err := rows.Scan(
			&a.Id,
			&a.Username,
			&a.Password,
			&a.NickName,
			&a.Phone,
			&a.Email,
			&a.Status,
			&a.RegisterDatetime,
			&a.CreateTimestamp,
			&a.ModifyTimestamp,
		); err != nil {
			return nil, err
		}
		data = append(data, &a)
	}
	return data, nil
}

func AccountWithUpdateOne(exec mysql.Exec, id string, values map[string]interface{}, timeout time.Duration) (int64, error) {
	if values == nil {
		return 0, fmt.Errorf("invalid parameter, value is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fields := []string{}
	args := make([]interface{}, 0, 6)
	for k, v := range values {
		if k == FieldAccountRegisterDatetime {
			continue
		}
		fields = append(fields, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	fields = append(fields, fmt.Sprintf("modify_timestamp = now()"))
	args = append(args, id)
	_sql := fmt.Sprintf(`update account set %s where id = ?`, strings.Join(fields, ", "))
	result, err := exec.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func AccountWithDeleteOne(exec mysql.Exec, id string, timeout time.Duration) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_sql := `delete from account where id = ?`

	result, err := exec.ExecContext(ctx, _sql, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func AccountWithInsertOne(exec mysql.Exec, account *Account, timeout time.Duration) (int64, error) {
	if account == nil {
		return 0, fmt.Errorf("invalid parameter, account is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fields := `(?, ?, ?, ?, ?, ?, ?, ?, now())`
	args := []interface{}{
		account.Id,
		account.Username,
		account.Password,
		account.NickName,
		account.Phone,
		account.Email,
		account.Status,
		account.RegisterDatetime,
	}
	_sql := fmt.Sprintf(`insert into account (%s) values %s`, strings.Join(AccountFields, ","), fields)
	result, err := exec.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

type Account struct {
	Id               string         `json:"id"`
	Username         sql.NullString `json:"username"`
	Password         string         `json:"password"`
	NickName         sql.NullString `json:"nick_name"`
	Phone            sql.NullString `json:"phone"`
	Email            sql.NullString `json:"email"`
	Status           int32          `json:"status"`
	RegisterDatetime time.Time      `json:"register_datetime"`
	CreateTimestamp  time.Time      `json:"create_timestamp"`
	ModifyTimestamp  sql.NullTime   `json:"modify_timestamp"`
}

const (
	FieldAccountId               = "id"
	FieldAccountUsername         = "username"
	FieldAccountPassword         = "password"
	FieldAccountNickName         = "nick_name"
	FieldAccountPhone            = "phone"
	FieldAccountEmail            = "email"
	FieldAccountStatus           = "status"
	FieldAccountRegisterDatetime = "register_datetime"
	FieldAccountCreateTimestamp  = "create_timestamp"
	FieldAccountModifyTimestamp  = "modify_timestamp"
)

var AccountFields = []string{
	FieldAccountId,
	FieldAccountUsername,
	FieldAccountPassword,
	FieldAccountNickName,
	FieldAccountPhone,
	FieldAccountEmail,
	FieldAccountStatus,
	FieldAccountRegisterDatetime,
	FieldAccountCreateTimestamp,
}
