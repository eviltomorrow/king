package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-auth/domain/persistence"
	"github.com/eviltomorrow/king/lib/auth"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/encrypt"
	"github.com/eviltomorrow/king/lib/snowflake"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

type (
	PassportStatus     int8
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

var (
	DefaultAccessTokenExpiresIn  = 5 * time.Minute
	DefaultRefreshTokenExpiresIn = 180 * time.Minute
)

type Passport struct {
	Id              string    `json:"id"`
	Account         string    `json:"account"`
	Code            string    `json:"code"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Status          int8      `json:"status"`
	CreateTimestamp time.Time `json:"create_timestamp"`
}

type Token struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	ExpiresIn    int64

	AccountId string
}

func PassportWithExist(ctx context.Context, account string) (bool, error) {
	if account == "" {
		return false, fmt.Errorf("account is nil")
	}

	_, err := persistence.PassportWithSelectOneByAccount(ctx, mysql.DB, account)
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
	data := &persistence.Passport{
		Id:           id,
		Account:      account,
		Salt:         s,
		SaltPassword: p,
		Status:       int8(NORMAL),
	}

	if _, err := persistence.PassportWithInsertOne(ctx, mysql.DB, data); err != nil {
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

	p, err := persistence.PassportWithSelectOneByAccount(ctx, mysql.DB, account)
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

func PassportWithChangeStatus(ctx context.Context, status PassportStatus, id string) error {
	_, err := persistence.PassportWithUpdateStatus(ctx, mysql.DB, int8(status), id)
	return err
}

func PassportWithChangePassword(ctx context.Context, password string, id string) error {
	if password == "" {
		return fmt.Errorf("accountpassword is nil")
	}

	s := encrypt.Salt()
	p := encrypt.Key(s, password)
	_, err := persistence.PassportWithUpdatePassword(ctx, mysql.DB, s, p, id)
	return err
}

func PassportWithRemove(ctx context.Context, id string) error {
	_, err := persistence.PassportWithDeleteOne(ctx, mysql.DB, id)
	return err
}

func PassportWithGet(ctx context.Context, account string) (*Passport, error) {
	p, err := persistence.PassportWithSelectOneByAccount(ctx, mysql.DB, account)
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

func PassportWithApplyToken(ctx context.Context, accountId, role string, accessTokenExpiresIn, refreshTokenExpiresIn time.Duration) (Token, error) {
	if accountId == "" {
		return Token{}, fmt.Errorf("accountId is nil")
	}

	accessToken, err := auth.JwtWithCreateToken(accountId, role, accessTokenExpiresIn)
	if err != nil {
		return Token{}, fmt.Errorf("create access_token failure, nest error: %v", err)
	}
	refreshToken, err := auth.JwtWithCreateToken(accountId, role, refreshTokenExpiresIn)
	if err != nil {
		return Token{}, fmt.Errorf("create refresh_token failure, nest error: %v", err)
	}

	token := Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenExpiresIn.Seconds()),
	}
	stateRefreshToken, err := auth.SwitchJwtTokenToStateToken(refreshToken)
	if err != nil {
		zlog.Warn("switch refresh_token to state_token failure failure", zap.Error(err), zap.String("accountId", accountId))
	} else {
		if err := auth.RenewStateToken(ctx, "", stateRefreshToken, accountId, refreshTokenExpiresIn); err != nil {
			zlog.Warn("RenewStateToken failure", zap.Error(err), zap.String("accountId", accountId))
		}
	}
	return token, nil
}

func PassportWithRenewToken(ctx context.Context, token Token) (Token, error) {
	if token.RefreshToken == "" {
		return Token{}, fmt.Errorf("refresh_token is nil")
	}

	cliams, err := auth.JwtWithVerifyToken(token.RefreshToken, nil)
	if err != nil {
		return Token{}, fmt.Errorf("verify refresh_token failure, nest error: %v", err)
	}

	stateRefreshToken, err := auth.SwitchJwtTokenToStateToken(token.RefreshToken)
	if err != nil {
		return Token{}, fmt.Errorf("switch refresh_token to state_token failure, nest error: %v", err)
	}
	ok, err := auth.SearchStateToken(ctx, stateRefreshToken)
	if err != nil {
		return Token{}, fmt.Errorf("search state_token failure, nest error: %v", err)
	}
	if !ok {
		return Token{}, fmt.Errorf("state_token is not found")
	}

	if err := auth.RevokeStateToken(ctx, stateRefreshToken); err != nil {
		zlog.Warn("revoke state_token failure", zap.Error(err), zap.String("token", stateRefreshToken))
	}
	return PassportWithApplyToken(ctx, cliams.AccountId, cliams.Role, DefaultAccessTokenExpiresIn, DefaultRefreshTokenExpiresIn)
}

func PassportWithVerifyToken(ctx context.Context, token Token) error {
	if token.AccessToken == "" {
		return fmt.Errorf("access_token is nil")
	}

	if _, err := auth.JwtWithVerifyToken(token.RefreshToken, nil); err != nil {
		return fmt.Errorf("verify access_token failure, nest error: %v", err)
	}
	return nil
}

func PassportWithRevokeToken(ctx context.Context, token Token) error {
	return nil
}
