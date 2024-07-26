package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/redis"
)

var (
	tokenPrefix        = "token_"
	tokenAccountPrefix = "token_account_"

	TokenLimitPerAccount int64 = 10
)

var ErrorNoAuth = errors.New("no auth")

func StateTokenWithParseJwtToken(jwtToken string) (string, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(jwtToken)); err != nil {
		return "", fmt.Errorf("panic: write sha256 failure, nest error: %v", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func stateTokenWithCleanExpires(ctx context.Context, id string) error {
	key := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	c := redis.Client.HGetAll(ctx, key)
	if err := c.Err(); err != nil {
		return err
	}
	data, err := c.Result()
	if err != nil {
		return err
	}

	exists := make([]string, 0, len(data))
	for k := range data {
		i := redis.Client.Exists(ctx, fmt.Sprintf("%s%s", tokenPrefix, k))
		if err := i.Err(); err != nil {
			return err
		}
		if v, err := i.Result(); err != nil {
			return err
		} else {
			if v == 0 {
				exists = append(exists, k)
			}
		}
	}
	if len(exists) != 0 {
		d := redis.Client.HDel(ctx, key, exists...)
		if err := d.Err(); err != nil {
			return err
		}
	}

	return nil
}

/*
TODO:

	LUA 改造
*/
func StateTokenWithCount(ctx context.Context, id string) (int64, error) {
	if err := stateTokenWithCleanExpires(ctx, id); err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	c := redis.Client.HLen(ctx, key)
	if err := c.Err(); err != nil {
		return 0, err
	}

	return c.Result()
}

func StateTokenWithSearchList(ctx context.Context, id string) ([]string, error) {
	if err := stateTokenWithCleanExpires(ctx, id); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	c := redis.Client.HGetAll(ctx, key)
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

	if oldToken != "" {
		key := fmt.Sprintf("%s%s", tokenPrefix, oldToken)

		c := redis.Client.Exists(ctx, key)
		if err := c.Err(); err != nil {
			return err
		}
		ok, err := c.Result()
		if err != nil {
			return err
		}
		if ok != 1 {
			return ErrorNoAuth
		}
	}

	count, err := StateTokenWithCount(ctx, id)
	if err != nil {
		return err
	}
	if count >= TokenLimitPerAccount {
		return fmt.Errorf("token apply has reached the maximum")
	}

	tokenKey := fmt.Sprintf("%s%s", tokenPrefix, newToken)
	s := redis.Client.Set(ctx, tokenKey, id, expiresIn)
	if err := s.Err(); err != nil {
		return err
	}

	accountKey := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	i := redis.Client.HSet(ctx, accountKey, newToken, 0)
	if err := i.Err(); err != nil {
		return err
	}

	if oldToken != "" {
		key1 := fmt.Sprintf("%s%s", tokenPrefix, oldToken)
		c := redis.Client.Del(ctx, key1)
		if err := c.Err(); err != nil {
			return err
		}

		key2 := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
		i := redis.Client.HDel(ctx, key2, key1)
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
	i := redis.Client.Del(ctx, keys...)
	if err := i.Err(); err != nil {
		return err
	}
	key := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	i = redis.Client.Del(ctx, key)
	if err := i.Err(); err != nil {
		return err
	}
	return nil
}

func StateTokenWithRevoke(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}

	tokenKey := fmt.Sprintf("%s%s", tokenPrefix, token)
	g := redis.Client.Get(ctx, tokenKey)
	if err := g.Err(); err != nil {
		return err
	}

	c := redis.Client.Del(ctx, tokenKey)
	if err := c.Err(); err != nil {
		return err
	}

	id, err := g.Result()
	if err != nil {
		return err
	}
	accountKey := fmt.Sprintf("%s%s", tokenAccountPrefix, id)
	i := redis.Client.HDel(ctx, accountKey, tokenKey)
	if err := i.Err(); err != nil {
		return err
	}
	return nil
}
