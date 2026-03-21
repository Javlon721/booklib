package userhandler

import (
	"fmt"
	"net/mail"
	"time"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/Javlon721/booklib/internal/types/password"
)

type User struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

func toHandlerUser(u userbus.User) User {
	return User{
		ID:          u.ID.String(),
		FirstName:   u.FirstName.String(),
		LastName:    u.LastName.String(),
		Email:       u.Email.Address,
		DateCreated: u.DateCreated.Format(time.RFC3339),
		DateUpdated: u.DateUpdated.Format(time.RFC3339),
	}
}

type NewUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func toBusUser(u NewUser) (userbus.CreateUser, error) {
	var errors errs.FieldErrors

	firstName, err := name.Parse(u.FirstName)

	if err != nil {
		errors.Add("first_name", err)
	}

	LastName, err := name.Parse(u.LastName)

	if err != nil {
		errors.Add("last_name", err)
	}

	password, err := password.Parse(u.Password)

	if err != nil {
		errors.Add("password", err)
	}

	email, err := mail.ParseAddress(u.Email)

	if err != nil {
		errors.Add("email", err)
	}

	if len(errors) > 0 {
		return userbus.CreateUser{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	return userbus.CreateUser{
		FirstName: firstName,
		LastName:  LastName,
		Email:     *email,
		Password:  password,
	}, nil
}
