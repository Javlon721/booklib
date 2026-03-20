package authhandler

import (
	"github.com/google/uuid"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateResp struct {
	UserID uuid.UUID
}
