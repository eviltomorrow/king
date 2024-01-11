package db

import (
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

func Init() {
	TruncateAccount()
}

func TestAccountWithInsertOne(t *testing.T) {
	_assert := assert.New(t)
	Init()

	for _, account := range []*Account{account1, account2, account3} {
		id := snowflake.GenerateID()
		account.Id = id
		affected, err := AccountWithInsertOne(mysql.DB, account, 10*time.Second)
		_assert.Nil(err)
		_assert.Equal(affected, int64(1))

		accountInserted, err := AccountWithSelectOne(mysql.DB, id, 10*time.Second)
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
	affected, err := AccountWithInsertOne(mysql.DB, account3, 10*time.Second)
	_assert.NotNil(err)
	_assert.Equal(int64(0), affected)

	affected, err = AccountWithInsertOne(mysql.DB, nil, 10*time.Second)
	_assert.NotNil(err)
	_assert.Equal(int64(0), affected)
}

func TestAccountWithDeleteOne(t *testing.T) {
	_assert := assert.New(t)
	Init()

	for _, account := range []*Account{account1, account2, account3} {
		id := snowflake.GenerateID()
		account.Id = id
		affected, err := AccountWithInsertOne(mysql.DB, account, 10*time.Second)
		_assert.Nil(err)
		_assert.Equal(int64(1), affected)

		affected, err = AccountWithDeleteOne(mysql.DB, id, 10*time.Second)
		_assert.Nil(err)
		_assert.Equal(int64(1), affected)

		_, err = AccountWithSelectOne(mysql.DB, id, 10*time.Second)
		_assert.NotNil(err)
		_assert.Equal(sql.ErrNoRows, err)
	}

	id := snowflake.GenerateID()
	affected, err := AccountWithDeleteOne(mysql.DB, id, 10*time.Second)
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)
}

func TestAccountWithUpdateOne(t *testing.T) {
	_assert := assert.New(t)
	Init()

	id := snowflake.GenerateID()
	account1.Id = id
	affected, err := AccountWithInsertOne(mysql.DB, account1, 10*time.Second)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	account1.Status = 2
	account1.Phone = sql.NullString{}
	account1.Password = "password1"
	values := map[string]interface{}{
		"status":   account1.Status,
		"phone":    account1.Phone,
		"password": account1.Password,
	}
	affected, err = AccountWithUpdateOne(mysql.DB, id, values, 10*time.Second)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	accountUpdated, err := AccountWithSelectOne(mysql.DB, id, 10*time.Second)
	_assert.Nil(err)
	_assert.Equal(int32(2), accountUpdated.Status)
	_assert.Equal(sql.NullString{}, accountUpdated.Phone)

	affected, err = AccountWithUpdateOne(mysql.DB, id, nil, 10*time.Second)
	_assert.NotNil(err)
	_assert.Equal(int64(0), affected)
}

func TestAccountWithSelectRange(t *testing.T) {
	_assert := assert.New(t)
	Init()

	accounts, err := AccountWithSelectRange(mysql.DB, 0, 10, 10*time.Second)
	_assert.Nil(err)
	_assert.Equal(0, len(accounts))
}
