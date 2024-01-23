package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPassportWithRegister(t *testing.T) {
	_assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := PassportWithRegister(ctx, "shepard", "liarsa")
	_assert.Nil(err)
	t.Logf("id: %s\r\n", id)
}

func TestPassportWithAuth(t *testing.T) {
	_assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := PassportWithAuth(ctx, PASSWORD, "shepard2", "liarsa")
	_assert.Nil(err)
}
