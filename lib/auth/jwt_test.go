package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJwtWithCreateToken(t *testing.T) {
	_assert := assert.New(t)

	token, err := JwtWithCreateToken("shepard", "admin", 3*time.Second)
	_assert.Nil(err)
	t.Logf("%s\r\n", token)

	// time.Sleep(10 * time.Second)
	token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoic2hlcGFyZCIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcwNjA3NDM5MiwibmJmIjoxNzA2MDc0Mzg5LCJpYXQiOjE3MDYwNzQzODl9.6gazgKo4zy2Y90vf6K1oxaRvS7Qcdo5U_DCYaFotbrI`
	c, err := JwtWithParseToken(token, nil)
	_assert.Nil(err)
	t.Logf("%s\r\n", c)
}
