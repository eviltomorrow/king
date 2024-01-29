package service

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/eviltomorrow/king/lib/db/redis"
	"github.com/stretchr/testify/assert"
)

func init() {
	redis.DSN = "redis://:admin123@localhost:6379/0?protocol=3"
	if err := redis.Connect(); err != nil {
		log.Fatal(err)
	}
}

var (
	id   = "shepard"
	role = "admin"
)

func TestTokenWithApply(t *testing.T) {
	_assert := assert.New(t)

	token, err := TokenWithApply(context.Background(), id, role)
	_assert.Nil(err)

	t.Log(token.RefreshToken, err)
	time.Sleep(2 * time.Second)

	token2, err := TokenWithRenew(context.Background(), token.RefreshToken)
	_assert.Nil(err)
	t.Log(token2.RefreshToken, id, err)
}
