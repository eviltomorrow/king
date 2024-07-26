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

type Token struct {
	AccessToken  string
	TokenType    string
	RefreshToken string

	ExpiresIn int64
}

func TokenWithApply(ctx context.Context, accountId, role string) (*Token, error) {
	if accountId == "" || role == "" {
		return nil, fmt.Errorf("accountId/role is nil")
	}

	accessToken, refreshToken, err := tokenWithGenerateTwo(accountId, role)
	if err != nil {
		return nil, err
	}

	stateToken, err := auth.StateTokenWithParseJwtToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if err := auth.StateTokenWithRenew(ctx, "", stateToken, accountId, RefreshTokenExpiresIn); err != nil {
		return nil, err
	}

	token := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(AccessTokenExpiresIn.Seconds()),
	}

	return token, nil
}

func TokenWithRenew(ctx context.Context, token string) (*Token, error) {
	if token == "" {
		return nil, fmt.Errorf("refresh_token is nil")
	}

	stateToken, err := auth.StateTokenWithParseJwtToken(token)
	if err != nil {
		return nil, err
	}

	claims, err := auth.JwtWithParseToken(token, nil)
	if err != nil {
		return nil, fmt.Errorf("parse refresh_token failure, nest error: %v", err)
	}

	accessToken, refreshToken, err := tokenWithGenerateTwo(claims.AccountId, claims.Role)
	if err != nil {
		return nil, err
	}

	stateRfreshToken, err := auth.StateTokenWithParseJwtToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if err := auth.StateTokenWithRenew(ctx, stateToken, stateRfreshToken, claims.AccountId, RefreshTokenExpiresIn); err != nil {
		return nil, err
	}

	data := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(AccessTokenExpiresIn.Seconds()),
	}

	return data, nil
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

func tokenWithGenerateTwo(accountId, role string) (string, string, error) {
	accessToken, err := auth.JwtWithCreateToken(accountId, role, AccessTokenExpiresIn)
	if err != nil {
		return "", "", fmt.Errorf("create access_token failure, nest error: %v", err)
	}
	refreshToken, err := auth.JwtWithCreateToken(accountId, role, RefreshTokenExpiresIn)
	if err != nil {
		return "", "", fmt.Errorf("create refresh_token failure, nest error: %v", err)
	}
	return accessToken, refreshToken, nil
}
