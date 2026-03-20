package authhandler

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	authbus "github.com/Javlon721/booklib/internal/bus/authBus"
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
	var parseErrors []string

	password, err := password.Parse(p.Password)

	if err != nil {
		parseErrors = append(parseErrors, fmt.Sprintf("password: %s", err.Error()))
	}

	email, err := mail.ParseAddress(p.Email)

	if err != nil {
		parseErrors = append(parseErrors, fmt.Sprintf("email: %s", err.Error()))
	}

	if len(parseErrors) > 0 {
		return authbus.AuthParams{}, errors.New(strings.Join(parseErrors, "\n"))
	}

	return authbus.AuthParams{
		Email:    *email,
		Password: password,
	}, nil
}
