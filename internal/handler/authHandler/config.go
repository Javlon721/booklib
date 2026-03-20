package authhandler

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	Secret         []byte
	TokenExpiresAt time.Duration
	Method         jwt.SigningMethod
}
