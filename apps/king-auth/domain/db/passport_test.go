package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/stretchr/testify/assert"
)

func TruncatePassport() error {
	_, err := mysql.DB.Exec("truncate table passport")
	return err
}

func InitPassport() {
	TruncatePassport()
}

var passport1 = Passport{
	Id:           "1",
	Account:      "liarsa",
	Code:         sql.NullString{},
	Salt:         "salt",
	SaltPassword: "salt_password",
	Status:       0,
	Email:        sql.NullString{String: "liarsa@qq.com", Valid: true},
	Phone:        sql.NullString{String: "13800138000", Valid: true},
}

func TestPassportWithInsertOne(t *testing.T) {
	_assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	InitPassport()

	for _, passport := range []Passport{passport1} {
		affected, err := PassportWithInsertOne(ctx, mysql.DB, &passport)
		_assert.Nil(err)
		_assert.Equal(affected, int64(1))

		passportInserted, err := PassportWithSelectOneById(ctx, mysql.DB, passport.Id)
		_assert.Nil(err)
		_assert.Equal(passportInserted.Id, passport.Id)
		_assert.Equal(passportInserted.Account, passport.Account)
		_assert.Equal(passportInserted.Code, passport.Code)
		_assert.Equal(passportInserted.Salt, passport.Salt)
		_assert.Equal(passportInserted.SaltPassword, passport.SaltPassword)
		_assert.Equal(passportInserted.Status, passport.Status)
		_assert.Equal(passportInserted.Email, passport.Email)
		_assert.Equal(passportInserted.Phone, passport.Phone)
	}
}
