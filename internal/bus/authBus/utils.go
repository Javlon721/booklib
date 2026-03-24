package authbus

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims jwt.Claims, secret []byte, method jwt.SigningMethod) (string, error) {
	token := jwt.NewWithClaims(method, claims)

	s, err := token.SignedString(secret)

	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return s, nil
}
