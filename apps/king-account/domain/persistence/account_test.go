package persistence

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/snowflake"
	"github.com/stretchr/testify/assert"
)

var account1 = &Account{
	Id:               "",
	Username:         sql.NullString{},
	Password:         "password",
	NickName:         sql.NullString{},
	Phone:            sql.NullString{},
	Email:            sql.NullString{},
	Status:           0,
	RegisterDatetime: time.Date(2024, time.January, 10, 10, 0o0, 0o0, 0o0, time.Local),
}

var account2 = &Account{
	Id:               "",
	Username:         sql.NullString{String: "shepard", Valid: true},
	Password:         "password",
	NickName:         sql.NullString{String: "Liarsa Shepard", Valid: true},
	Phone:            sql.NullString{String: "13800138000", Valid: true},
	Email:            sql.NullString{String: "shepard@qq.com", Valid: true},
	Status:           1,
	RegisterDatetime: time.Date(2024, time.January, 10, 12, 0o0, 0o0, 0o0, time.Local),
}

var account3 = &Account{
	Id:               "",
	Username:         sql.NullString{String: "liarsa", Valid: true},
	Password:         "password",
	NickName:         sql.NullString{String: "Liarsa Shepard", Valid: true},
	Phone:            sql.NullString{},
	Email:            sql.NullString{String: "liarsa@qq.com", Valid: true},
	Status:           1,
	RegisterDatetime: time.Date(2024, time.January, 10, 12, 0o0, 0o0, 0o0, time.Local),
}

func TruncateAccount() error {
	_, err := mysql.DB.Exec("truncate table account")
	return err
}

func InitAccount() {
	TruncateAccount()
}

func TestAccountWithInsertOne(t *testing.T) {
	_assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	InitAccount()

	for _, account := range []*Account{account1, account2, account3} {
		account.Id = snowflake.GenerateID()
		affected, err := AccountWithInsertOne(ctx, mysql.DB, account)
		_assert.Nil(err)
		_assert.Equal(affected, int64(1))

		accountInserted, err := AccountWithSelectOne(ctx, mysql.DB, account.Id)
		_assert.Nil(err)
		_assert.Equal(accountInserted.Id, account.Id)
		_assert.Equal(accountInserted.Username, account.Username)
		_assert.Equal(accountInserted.Password, account.Password)
		_assert.Equal(accountInserted.NickName, account.NickName)
		_assert.Equal(accountInserted.Phone, account.Phone)
		_assert.Equal(accountInserted.Email, account.Email)
		_assert.Equal(accountInserted.Status, account.Status)
		_assert.Equal(accountInserted.RegisterDatetime, account.RegisterDatetime)
	}

	account3.Id = snowflake.GenerateID()
	affected, err := AccountWithInsertOne(ctx, mysql.DB, account3)
	_assert.NotNil(err)
	_assert.Equal(int64(0), affected)

	affected, err = AccountWithInsertOne(ctx, mysql.DB, nil)
	_assert.NotNil(err)
	_assert.Equal(int64(0), affected)
}

func TestAccountWithDeleteOne(t *testing.T) {
	_assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	InitAccount()

	for _, account := range []*Account{account1, account2, account3} {
		id := snowflake.GenerateID()
		account.Id = id
		affected, err := AccountWithInsertOne(ctx, mysql.DB, account)
		_assert.Nil(err)
		_assert.Equal(int64(1), affected)

		affected, err = AccountWithDeleteOne(ctx, mysql.DB, id)
		_assert.Nil(err)
		_assert.Equal(int64(1), affected)

		_, err = AccountWithSelectOne(ctx, mysql.DB, id)
		_assert.NotNil(err)
		_assert.Equal(sql.ErrNoRows, err)
	}

	id := snowflake.GenerateID()
	affected, err := AccountWithDeleteOne(ctx, mysql.DB, id)
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)
}

func TestAccountWithUpdateOne(t *testing.T) {
	_assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	InitAccount()

	id := snowflake.GenerateID()
	account1.Id = id
	affected, err := AccountWithInsertOne(ctx, mysql.DB, account1)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	account1.Status = 2
	account1.Phone = sql.NullString{}
	account1.Password = "password1"
	affected, err = AccountWithUpdateOne(ctx, mysql.DB, account1, id)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	accountUpdated, err := AccountWithSelectOne(ctx, mysql.DB, id)
	_assert.Nil(err)
	_assert.Equal(int8(2), accountUpdated.Status)
	_assert.Equal(sql.NullString{}, accountUpdated.Phone)

	affected, err = AccountWithUpdateOne(ctx, mysql.DB, nil, id)
	_assert.NotNil(err)
	_assert.Equal(int64(0), affected)
}

func TestAccountWithSelectRange(t *testing.T) {
	_assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	InitAccount()

	accounts, err := AccountWithSelectRange(ctx, mysql.DB, 0, 10)
	_assert.Nil(err)
	_assert.Equal(0, len(accounts))

	for _, account := range []*Account{account1, account2, account3} {
		id := snowflake.GenerateID()
		account.Id = id
		affected, err := AccountWithInsertOne(ctx, mysql.DB, account)
		_assert.Nil(err)
		_assert.Equal(int64(1), affected)
	}

	accounts, err = AccountWithSelectRange(ctx, mysql.DB, 0, 2)
	_assert.Nil(err)
	_assert.Equal(2, len(accounts))
}
