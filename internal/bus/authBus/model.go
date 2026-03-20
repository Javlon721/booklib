package authbus

import (
	"net/mail"

	"github.com/Javlon721/booklib/internal/types/password"
	"github.com/google/uuid"
)

type AuthenticateResp struct {
	UserID uuid.UUID
}

type AuthParams struct {
	Email    mail.Address
	Password password.Password
}
