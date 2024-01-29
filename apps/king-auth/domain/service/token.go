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

func TokenWithTryApply(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("id is nil")
	}
	return false, nil
}

func TokenWithApply(ctx context.Context, id, role string, expired *Token) (Token, string, error) {
	if id == "" || role == "" {
		return Token{}, "", fmt.Errorf("id/role is nil")
	}

	accessToken, err := auth.JwtWithCreateToken(id, role, AccessTokenExpiresIn)
	if err != nil {
		return Token{}, "", fmt.Errorf("create access_token failure, nest error: %v", err)
	}
	refreshToken, err := auth.JwtWithCreateToken(id, role, RefreshTokenExpiresIn)
	if err != nil {
		return Token{}, "", fmt.Errorf("create refresh_token failure, nest error: %v", err)
	}

	token := Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(AccessTokenExpiresIn.Seconds()),
	}

	stateRefreshToken, err := auth.StateTokenWithParseJwtToken(token.RefreshToken)
	if err != nil {
		return token, id, err
	} else {
		var expiredToken string
		if expired != nil {
			expiredToken = expired.RefreshToken
		}
		if err := auth.StateTokenWithRenew(ctx, expiredToken, stateRefreshToken, id, RefreshTokenExpiresIn); err != nil {
			return token, id, err
		}
	}

	return token, id, nil
}

func TokenWithRenew(ctx context.Context, expired Token) (Token, string, error) {
	if expired.RefreshToken == "" {
		return Token{}, "", fmt.Errorf("refresh_token is nil")
	}

	cliams, err := auth.JwtWithParseToken(expired.RefreshToken, nil)
	if err != nil {
		return Token{}, "", fmt.Errorf("parse refresh_token failure, nest error: %v", err)
	}

	return TokenWithApply(ctx, cliams.AccountId, cliams.Role, &expired)
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
