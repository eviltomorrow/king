package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
)

func AccountWithModifyOne(exec mysql.Exec, account *Account, timeout time.Duration) (int64, error) {
	if account == nil {
		return 0, fmt.Errorf("invalid parameter, account is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fields := []string{}
	args := make([]interface{}, 0, 6)

	_sql := `update account set where id = ?`
	return 0, nil
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

type Value interface {
	Value() (driver.Value, error)
}

func FieldWithSQL(columnName string, value Value) (string, interface{}, bool) {
	if v, err := value.Value(); err == nil {
		return fmt.Sprintf("%s = ?", columnName), v, true
	}
	return "", nil, false
}

func FieldWithString(columnName string, value string) (string, interface{}, bool) {
	if value == "" {
		return "", nil, false
	}
	return fmt.Sprintf("%s = ?", columnName), value, true
}

func FieldWithNumber[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64](columnName string, value T, ignore bool) (string, interface{}, bool) {
	if value == 0 {
		return "", nil, false
	}
	return fmt.Sprintf("%s = ?", columnName), value, true
}
