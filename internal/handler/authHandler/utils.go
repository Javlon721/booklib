package authhandler

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims jwt.MapClaims, secret []byte, method jwt.SigningMethod) (string, error) {
	token := jwt.NewWithClaims(method, claims)

	s, err := token.SignedString(secret)

	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return s, nil
}

func ParseToken(tokenString string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	return token, err
}
