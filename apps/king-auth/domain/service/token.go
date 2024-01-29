package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/auth"
)

var (
	AccessTokenExpiresIn  = 10 * time.Minute
	RefreshTokenExpiresIn = 120 * time.Minute
)

func TokenWithApply(ctx context.Context, accountId, role string, expired *Token) (*Token, error) {
	if accountId == "" || role == "" {
		return nil, fmt.Errorf("accountId/role is nil")
	}

	accessToken, err := auth.JwtWithCreateToken(accountId, role, AccessTokenExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("create access_token failure, nest error: %v", err)
	}
	refreshToken, err := auth.JwtWithCreateToken(accountId, role, RefreshTokenExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("create refresh_token failure, nest error: %v", err)
	}

	token := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(AccessTokenExpiresIn.Seconds()),

		AccountId: accountId,
		Role:      role,
	}
	if err := token.Parse(); err != nil {
		return nil, fmt.Errorf("parse state_token failure, nest error: %v", err)
	}

	expiredStateToken := ""
	if expired != nil {
		expiredStateToken = expired.StateRefreshToken
	}
	if err := auth.StateTokenWithRenew(ctx, expiredStateToken, token.StateRefreshToken, accountId, RefreshTokenExpiresIn); err != nil {
		return token, err
	}

	return token, nil
}

func TokenWithRenew(ctx context.Context, expired *Token) (*Token, error) {
	if expired == nil || expired.RefreshToken == "" {
		return nil, fmt.Errorf("refresh_token is nil")
	}

	if err := expired.Parse(); err != nil {
		return nil, err
	}

	return TokenWithApply(ctx, expired.AccountId, expired.Role, expired)
}

func TokenWithVerify(ctx context.Context, token Token) error {
	if token.AccessToken == "" {
		return fmt.Errorf("access_token is nil")
	}

	if _, err := auth.JwtWithParseToken(token.RefreshToken, nil); err != nil {
		return fmt.Errorf("parse access_token failure, nest error: %v", err)
	}
	return nil
}

func TokenWithRevokeByToken(ctx context.Context, token Token) error {
	if token.RefreshToken == "" {
		return fmt.Errorf("refresh_token is nil")
	}

	stateRefreshToken, err := auth.StateTokenWithParseJwtToken(token.RefreshToken)
	if err != nil {
		return fmt.Errorf("switch refresh_token to state_token failure, nest error: %v", err)
	}

	return auth.StateTokenWithRevoke(ctx, stateRefreshToken)
}

func TokenWithRevokeByAccountId(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is nil")
	}

	return auth.StateTokenWithRevokeAll(ctx, id)
}

type Token struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	ExpiresIn    int64

	AccountId         string
	Role              string
	StateRefreshToken string
}

func (t *Token) Parse() error {
	if t == nil {
		return nil
	}

	if t.AccountId == "" || t.Role == "" {
		claims, err := auth.JwtWithParseToken(t.AccessToken, nil)
		if err != nil {
			return err
		}
		t.Role = claims.Role
		t.AccountId = claims.AccountId
	}

	if t.StateRefreshToken == "" {
		token, err := auth.StateTokenWithParseJwtToken(t.RefreshToken)
		if err != nil {
			return err
		}
		t.StateRefreshToken = token
	}

	return nil
}
