package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-account/domain/persistence"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/encrypt"
	"github.com/eviltomorrow/king/lib/snowflake"
)

const (
	NORMAL = iota
	LOCK
)

type Account struct {
	Id               string    `json:"id"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	NickName         string    `json:"nick_name"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email"`
	Status           int8      `json:"status"`
	RegisterDatetime time.Time `json:"register_datetime"`
}

func AccountWithRegister(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("invalid username or password")
	}

	id := snowflake.GenerateID()

	p, err := encrypt.Key(password)
	if err != nil {
		return "", err
	}
	data := &persistence.Account{
		Id:               id,
		Username:         sql.NullString{String: username, Valid: true},
		Password:         p,
		Status:           NORMAL,
		RegisterDatetime: time.Now(),
	}

	if _, err := persistence.AccountWithInsertOne(ctx, mysql.DB, data); err != nil {
		return "", err
	}
	return id, nil
}
