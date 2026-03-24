package authbus

import (
	"net/mail"

	"github.com/Javlon721/booklib/internal/types/password"
	"github.com/Javlon721/booklib/internal/types/role"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthenticateResp struct {
	UserID uuid.UUID
	Roles  []role.Role
}

type AuthParams struct {
	Email    mail.Address
	Password password.Password
}

type Claims struct {
	jwt.RegisteredClaims
	Roles []string
}
