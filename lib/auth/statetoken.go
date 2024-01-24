package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/redis"
)

var tokenPrefix = "token_"

func SwithJwtTokenToStateToken(jwtToken string) (string, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(jwtToken)); err != nil {
		return "", fmt.Errorf("panic: write sha256 failure, nest error: %v", err)
	}
	return fmt.Sprintf("%s%x", tokenPrefix, h.Sum(nil)), nil
}

func SearchStateToken(ctx context.Context, token string) (bool, error) {
	c := redis.RDB.Get(ctx, token)
	if err := c.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func RevokeStateToken(ctx context.Context, token string) error {
	c := redis.RDB.Del(ctx, token)
	if err := c.Err(); err != nil {
		return err
	}
	return nil
}

func RenewStateToken(ctx context.Context, oldToken, newToken string, account string, expiresIn time.Duration) error {
	if oldToken != "" {
		c := redis.RDB.Del(ctx, oldToken)
		if err := c.Err(); err != nil {
			return err
		}
	}

	s := redis.RDB.Set(ctx, newToken, account, expiresIn)
	if err := s.Err(); err != nil {
		return err
	}
	return nil
}
