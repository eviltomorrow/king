package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/redis"
)

var (
	tokenPrefix        = "token_"
	tokenAccountPrefix = "token_account_"

	TokenLimitPerAccount int64 = 10
)

func StateTokenWithParseJwtToken(jwtToken string) (string, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(jwtToken)); err != nil {
		return "", fmt.Errorf("panic: write sha256 failure, nest error: %v", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func StateTokenWithCount(ctx context.Context, id string) (int64, error) {
	key := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	c := redis.RDB.HLen(ctx, key)
	if err := c.Err(); err != nil {
		return 0, err
	}
	return c.Result()
}

func StateTokenWithSearch(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("%s%s", tokenPrefix, token)
	c := redis.RDB.Get(ctx, key)
	if err := c.Err(); err != nil {
		return "", err
	}
	return c.Result()
}

func StateTokenWithSearchList(ctx context.Context, id string) ([]string, error) {
	key := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	c := redis.RDB.HGetAll(ctx, key)
	if err := c.Err(); err != nil {
		return nil, err
	}
	data, err := c.Result()
	if err != nil {
		return nil, err
	}

	tokens := make([]string, 0, len(data))
	for k := range data {
		tokens = append(tokens, k)
	}

	return tokens, nil
}

func StateTokenWithRenew(ctx context.Context, oldToken, newToken string, id string, expiresIn time.Duration) error {
	if newToken == "" || id == "" {
		return fmt.Errorf("new_token/id is nil")
	}

	count, err := StateTokenWithCount(ctx, id)
	if err != nil {
		return err
	}
	if count >= TokenLimitPerAccount {
		return fmt.Errorf("token apply has reached the maximum")
	}

	tokenKey := fmt.Sprintf("%s%s", tokenPrefix, newToken)
	s := redis.RDB.Set(ctx, tokenKey, id, expiresIn)
	if err := s.Err(); err != nil {
		return err
	}

	accountKey := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	i := redis.RDB.HSet(ctx, accountKey, newToken, 0)
	if err := i.Err(); err != nil {
		return err
	}
	b := redis.RDB.Expire(ctx, fmt.Sprintf("%s:%s", accountKey, newToken), expiresIn)
	if err := b.Err(); err != nil {
		return err
	}

	if oldToken != "" {
		key1 := fmt.Sprintf("%s%s", tokenPrefix, oldToken)
		c := redis.RDB.Del(ctx, key1)
		if err := c.Err(); err != nil {
			return err
		}

		key2 := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
		i := redis.RDB.HDel(ctx, key2, key1)
		if err := i.Err(); err != nil {
			return err
		}
	}

	return nil
}

func StateTokenWithRevokeAll(ctx context.Context, id string) error {
	tokens, err := StateTokenWithSearchList(ctx, id)
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(tokens))
	for _, token := range tokens {
		key := fmt.Sprintf("%s%s", tokenPrefix, token)
		keys = append(keys, key)
	}
	i := redis.RDB.Del(ctx, keys...)
	if err := i.Err(); err != nil {
		return err
	}
	key := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	i = redis.RDB.Del(ctx, key)
	if err := i.Err(); err != nil {
		return err
	}
	return nil
}

func StateTokenWithRevoke(ctx context.Context, token string) error {
	id, err := StateTokenWithSearch(ctx, token)
	if err != nil {
		return err
	}
	if id == "" {
		return nil
	}

	tokenKey := fmt.Sprintf("%s%s", tokenPrefix, token)
	c := redis.RDB.Del(ctx, tokenKey)
	if err := c.Err(); err != nil {
		return err
	}

	accountKey := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	i := redis.RDB.HDel(ctx, accountKey, tokenKey)
	if err := i.Err(); err != nil {
		return err
	}
	return nil
}
