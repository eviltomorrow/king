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
	DefaultAccessTokenExpiresIn  = 5 * time.Minute
	DefaultRefreshTokenExpiresIn = 120 * time.Minute
)

func TokenWithApply(ctx context.Context, accountId, role string, accessTokenExpiresIn, refreshTokenExpiresIn time.Duration) (Token, error) {
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
		zlog.Warn("switch refresh_token to state_token failure", zap.Error(err), zap.String("accountId", accountId))
	} else {
		if err := auth.RenewStateToken(ctx, "", stateRefreshToken, accountId, refreshTokenExpiresIn); err != nil {
			zlog.Warn("RenewStateToken failure", zap.Error(err), zap.String("accountId", accountId))
		}
	}
	return token, nil
}

func TokenWithRenew(ctx context.Context, token Token) (Token, error) {
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
	return TokenWithApply(ctx, cliams.AccountId, cliams.Role, DefaultAccessTokenExpiresIn, DefaultRefreshTokenExpiresIn)
}

func TokenWithVerify(ctx context.Context, token Token) error {
	if token.AccessToken == "" {
		return fmt.Errorf("access_token is nil")
	}

	if _, err := auth.JwtWithVerifyToken(token.RefreshToken, nil); err != nil {
		return fmt.Errorf("verify access_token failure, nest error: %v", err)
	}
	return nil
}

func TokenWithRevoke(ctx context.Context, token Token) error {
	if token.RefreshToken == "" {
		return fmt.Errorf("refresh_token is nil")
	}

	stateRefreshToken, err := auth.SwitchJwtTokenToStateToken(token.RefreshToken)
	if err != nil {
		return fmt.Errorf("switch refresh_token to state_token failure, nest error: %v", err)
	}

	return auth.RevokeStateToken(ctx, stateRefreshToken)
}
