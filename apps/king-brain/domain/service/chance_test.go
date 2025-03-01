package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindPossibleChance(t *testing.T) {
	assert := assert.New(t)

	current, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	assert.Nil(err)

	plans := FindPossibleChance(context.Background(), current)
	fmt.Println(len(plans))
}
