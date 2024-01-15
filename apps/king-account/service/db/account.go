package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/orm"
)

func AccountWithSelectOne(ctx context.Context, exec mysql.Exec, id string) (*Account, error) {
	account := Account{}
	scan := func(row *sql.Row) error {
		return row.Scan(
			&account.Id,
			&account.Username,
			&account.Password,
			&account.NickName,
			&account.Phone,
			&account.Email,
			&account.Status,
			&account.RegisterDatetime,
			&account.CreateTimestamp,
			&account.ModifyTimestamp,
		)
	}
	if err := orm.TableWithSelectOne(ctx, exec, TableAccountName, AccountFields, map[string]interface{}{FieldAccountId: id}, scan); err != nil {
		return nil, err
	}
	return &account, nil
}

func AccountWithSelectRange(ctx context.Context, exec mysql.Exec, offset, limit int64) ([]*Account, error) {
	accounts := make([]*Account, 0, limit)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			account := Account{}
			if err := rows.Scan(
				&account.Id,
				&account.Username,
				&account.Password,
				&account.NickName,
				&account.Phone,
				&account.Email,
				&account.Status,
				&account.RegisterDatetime,
				&account.CreateTimestamp,
				&account.ModifyTimestamp,
			); err != nil {
				return err
			}
			accounts = append(accounts, &account)
		}
		return nil
	}
	if err := orm.TableWithSelectRange(ctx, exec, TableAccountName, AccountFields, nil, nil, offset, limit, scan); err != nil {
		return nil, err
	}
	return accounts, nil
}

func AccountWithUpdateOne(ctx context.Context, exec mysql.Exec, id string, value map[string]interface{}) (int64, error) {
	if value == nil {
		return 0, fmt.Errorf("invalid parameter, value is nil")
	}
	return orm.TableWithUpdate(ctx, exec, TableAccountName, value, map[string]interface{}{FieldAccountId: id})
}

func AccountWithDeleteOne(ctx context.Context, exec mysql.Exec, id string) (int64, error) {
	if id == "" {
		return 0, fmt.Errorf("invalid parameter, id is nil")
	}
	return orm.TableWithDelete(ctx, exec, TableAccountName, map[string]interface{}{FieldAccountId: id})
}

func AccountWithInsertOne(ctx context.Context, exec mysql.Exec, account *Account) (int64, error) {
	if account == nil {
		return 0, fmt.Errorf("invalid parameter, account is nil")
	}

	value := map[string]interface{}{
		FieldAccountId:               account.Id,
		FieldAccountUsername:         account.Username,
		FieldAccountPassword:         account.Password,
		FieldAccountNickName:         account.NickName,
		FieldAccountPhone:            account.Phone,
		FieldAccountEmail:            account.Email,
		FieldAccountStatus:           account.Status,
		FieldAccountRegisterDatetime: account.RegisterDatetime,
	}
	return orm.TableWithInsert(ctx, exec, TableAccountName, value)
}

type Account struct {
	Id               string         `json:"id"`
	Username         sql.NullString `json:"username"`
	Password         string         `json:"password"`
	NickName         sql.NullString `json:"nick_name"`
	Phone            sql.NullString `json:"phone"`
	Email            sql.NullString `json:"email"`
	Status           int8           `json:"status"`
	RegisterDatetime time.Time      `json:"register_datetime"`
	CreateTimestamp  time.Time      `json:"create_timestamp"`
	ModifyTimestamp  sql.NullTime   `json:"modify_timestamp"`
}

const (
	TableAccountName = "account"

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
	FieldAccountModifyTimestamp,
}
