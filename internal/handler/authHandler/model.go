package authhandler

import (
	"fmt"
	"net/mail"

	authbus "github.com/Javlon721/booklib/internal/bus/authBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/types/password"
	"github.com/google/uuid"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateResp struct {
	UserID uuid.UUID
}

func toBusAuthParams(p AuthParams) (authbus.AuthParams, error) {
	var errors errs.FieldErrors

	password, err := password.Parse(p.Password)

	if err != nil {
		errors.Add("password", err)
	}

	email, err := mail.ParseAddress(p.Email)

	if err != nil {
		errors.Add("email", err)
	}

	if len(errors) > 0 {
		return authbus.AuthParams{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	return authbus.AuthParams{
		Email:    *email,
		Password: password,
	}, nil
}
