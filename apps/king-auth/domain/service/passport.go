package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-auth/domain/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/encrypt"
	"github.com/eviltomorrow/king/lib/snowflake"
)

type (
	PassportStatus     int32
	PassportAuthMethod int8
)

const (
	NORMAL PassportStatus = iota
	LOCK
)

const (
	PASSWORD PassportAuthMethod = iota
	SMS
)

var ErrPassportNoAccount = errors.New("no account")

type Passport struct {
	Id              string    `json:"id"`
	Account         string    `json:"account"`
	Code            string    `json:"code"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Status          int32     `json:"status"`
	CreateTimestamp time.Time `json:"create_timestamp"`
}

func PassportWithExist(ctx context.Context, account string) (bool, error) {
	if account == "" {
		return false, fmt.Errorf("account is nil")
	}

	_, err := db.PassportWithSelectOneByAccount(ctx, mysql.DB, account)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func PassportWithRegister(ctx context.Context, account, password string) (string, error) {
	if account == "" || password == "" {
		return "", fmt.Errorf("account/password is nil")
	}

	id := fmt.Sprintf("k_%s", snowflake.GenerateID())

	s := encrypt.Salt()
	p := encrypt.Key(s, password)
	data := &db.Passport{
		Id:           id,
		Account:      account,
		Salt:         s,
		SaltPassword: p,
		Status:       int32(NORMAL),
	}

	if _, err := db.PassportWithInsertOne(ctx, mysql.DB, data); err != nil {
		return "", err
	}
	return id, nil
}

func PassportWithAuth(ctx context.Context, method PassportAuthMethod, account, key string) (*Passport, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account/key is nil")
	}

	switch method {
	case PASSWORD:
		return authWithPassword(ctx, account, key)
	case SMS:
		return nil, fmt.Errorf("not implement")
	default:
		return nil, fmt.Errorf("panic: not support auth method")
	}
}

func authWithPassword(ctx context.Context, account, password string) (*Passport, error) {
	if account == "" || password == "" {
		return nil, fmt.Errorf("account/password is nil")
	}

	p, err := db.PassportWithSelectOneByAccount(ctx, mysql.DB, account)
	if err == sql.ErrNoRows {
		return nil, ErrPassportNoAccount
	}
	if err != nil {
		return nil, err
	}
	key := encrypt.Key(p.Salt, password)
	if p.SaltPassword != key {
		return nil, fmt.Errorf("incorrect account/password")
	}
	return &Passport{
		Id:              p.Id,
		Account:         p.Account,
		Code:            p.Code.String,
		Email:           p.Email.String,
		Phone:           p.Phone.String,
		Status:          p.Status,
		CreateTimestamp: p.CreateTimestamp,
	}, nil
}

func PassportWithChangeStatus(ctx context.Context, status PassportStatus, accountId string) error {
	if accountId == "" {
		return fmt.Errorf("accountId is nil")
	}
	switch status {
	case NORMAL:
	case LOCK:
	default:
		return fmt.Errorf("not support status[%v]", status)
	}
	_, err := db.PassportWithUpdateStatus(ctx, mysql.DB, int8(status), accountId)
	return err
}

func PassportWithChangePassword(ctx context.Context, password string, accountId string) error {
	if password == "" {
		return fmt.Errorf("password is nil")
	}
	if accountId == "" {
		return fmt.Errorf("accountId is nil")
	}

	s := encrypt.Salt()
	p := encrypt.Key(s, password)
	_, err := db.PassportWithUpdatePassword(ctx, mysql.DB, s, p, accountId)
	return err
}

func PassportWithRemove(ctx context.Context, accountId string) error {
	if accountId == "" {
		return fmt.Errorf("id is nil")
	}
	_, err := db.PassportWithDeleteOne(ctx, mysql.DB, accountId)
	return err
}

func PassportWithGet(ctx context.Context, account string) (*Passport, error) {
	if account == "" {
		return nil, fmt.Errorf("account is nil")
	}
	p, err := db.PassportWithSelectOneByAccount(ctx, mysql.DB, account)
	if err == sql.ErrNoRows {
		return nil, ErrPassportNoAccount
	}
	if err != nil {
		return nil, err
	}
	return &Passport{
		Id:              p.Id,
		Account:         p.Account,
		Code:            p.Code.String,
		Email:           p.Email.String,
		Phone:           p.Phone.String,
		Status:          p.Status,
		CreateTimestamp: p.CreateTimestamp,
	}, nil
}
