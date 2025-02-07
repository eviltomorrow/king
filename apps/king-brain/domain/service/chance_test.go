package service

import (
	"context"
	"testing"
	"time"

	_ "github.com/eviltomorrow/king/apps/king-brain/domain/feature"
	"github.com/stretchr/testify/assert"
)

func TestFindPossibleChance(t *testing.T) {
	assert := assert.New(t)

	current, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	assert.Nil(err)

	FindPossibleChance(context.Background(), current)
}
