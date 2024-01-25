package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/auth"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

var (
	DefaultAccessTokenExpiresIn  = 10 * time.Minute
	DefaultRefreshTokenExpiresIn = 720 * time.Minute
)

func TokenWithTryApply(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("id is nil")
	}
	return false, nil
}

func TokenWithApply(ctx context.Context, id, role string, accessTokenExpiresIn, refreshTokenExpiresIn time.Duration) (Token, string, error) {
	if id == "" || role == "" {
		return Token{}, "", fmt.Errorf("id/role is nil")
	}
	if accessTokenExpiresIn < DefaultAccessTokenExpiresIn || accessTokenExpiresIn > refreshTokenExpiresIn {
		return Token{}, "", fmt.Errorf("expires_in is wrong")
	}

	accessToken, err := auth.JwtWithCreateToken(id, role, accessTokenExpiresIn)
	if err != nil {
		return Token{}, "", fmt.Errorf("create access_token failure, nest error: %v", err)
	}
	refreshToken, err := auth.JwtWithCreateToken(id, role, refreshTokenExpiresIn)
	if err != nil {
		return Token{}, "", fmt.Errorf("create refresh_token failure, nest error: %v", err)
	}

	token := Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTokenExpiresIn.Seconds()),
	}

	stateRefreshToken, err := auth.StateTokenWithParseJwtToken(token.RefreshToken)
	if err != nil {
		zlog.Warn("StateTokenWithParseJwtToken failure", zap.Error(err), zap.String("accountId", id))
	} else {
		if err := auth.StateTokenWithRenew(ctx, "", stateRefreshToken, id, refreshTokenExpiresIn); err != nil {
			zlog.Warn("RenewStateToken failure", zap.Error(err), zap.String("accountId", id))
		}
	}

	return token, id, nil
}

func TokenWithRenew(ctx context.Context, token Token) (Token, string, error) {
	if token.RefreshToken == "" {
		return Token{}, "", fmt.Errorf("refresh_token is nil")
	}

	cliams, err := auth.JwtWithParseToken(token.RefreshToken, nil)
	if err != nil {
		return Token{}, "", fmt.Errorf("parse refresh_token failure, nest error: %v", err)
	}

	return TokenWithApply(ctx, cliams.AccountId, cliams.Role, DefaultAccessTokenExpiresIn, DefaultRefreshTokenExpiresIn)
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
