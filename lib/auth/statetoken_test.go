package auth

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/eviltomorrow/king/lib/db/redis"
	"github.com/stretchr/testify/assert"
)

func init() {
	redis.DSN = "redis://:admin123@localhost:6379/1?protocol=3"
	if err := redis.Connect(); err != nil {
		log.Fatal(err)
	}
}

func TestRenewStateToken(t *testing.T) {
	_assert := assert.New(t)

	expiresIn := 10 * time.Minute
	token, err := JwtWithCreateToken("shepard", "admin", expiresIn)
	_assert.Nil(err)
	t.Logf("%s\r\n", token)

	stateToken, err := SwithJwtTokenToStateToken(token)
	_assert.Nil(err)
	t.Logf("%s\r\n", stateToken)

	err = RenewStateToken(context.Background(), "", stateToken, "shepard", expiresIn)
	_assert.Nil(err)

	ok, err := SearchStateToken(context.Background(), stateToken)
	_assert.Nil(err)
	_assert.Equal(true, ok)
}
