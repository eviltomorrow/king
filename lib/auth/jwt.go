package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	SigningKey      = []byte("123")
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenIllegal = errors.New("token is illegal")
)

type JwtClaims struct {
	AccountId string `json:"account_id"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

func JwtWithCreateToken(accountId string, role string, expireIn time.Duration) (string, error) {
	claims := JwtClaims{
		AccountId: accountId,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(SigningKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func JwtWithParseToken(tokenStr string, f func(string) error) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})

	switch {
	case token.Valid:
		claims, ok := token.Claims.(*JwtClaims)
		if ok {
			return claims, nil
		}
		if f != nil {
			return claims, f(claims.AccountId)
		}
		return claims, nil

	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, ErrTokenExpired

	default:
		return nil, ErrTokenIllegal

	}
}
