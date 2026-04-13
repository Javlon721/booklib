package userhandler

import (
	"fmt"
	"net/mail"
	"time"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/Javlon721/booklib/internal/types/password"
	"github.com/Javlon721/booklib/internal/types/role"
)

type User struct {
	ID          string   `json:"id"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name,omitempty"`
	Email       string   `json:"email"`
	DateCreated string   `json:"date_created"`
	DateUpdated string   `json:"date_updated"`
	Roles       []string `json:"roles"`
}

func toHandlerUser(u userbus.User) User {
	return User{
		ID:          u.ID.String(),
		FirstName:   u.FirstName.String(),
		LastName:    u.LastName.String(),
		Email:       u.Email.Address,
		DateCreated: u.DateCreated.Format(time.RFC3339),
		DateUpdated: u.DateUpdated.Format(time.RFC3339),
		Roles:       role.ParseToString(u.Roles),
	}
}

type NewUser struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name,omitempty"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles"`
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

	roles, err := role.ParseMany(u.Roles)

	if err != nil {
		errors.Add("roles", err)
	}

	if len(errors) > 0 {
		return userbus.CreateUser{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	return userbus.CreateUser{
		FirstName: firstName,
		LastName:  LastName,
		Email:     *email,
		Password:  password,
		Roles:     roles,
	}, nil
}

// -------------------------------------------------------------------------
type UpdateUser struct {
	FirstName *string  `json:"first_name"`
	LastName  *string  `json:"last_name,omitempty"`
	Email     *string  `json:"email"`
	Password  *string  `json:"password"`
	Roles     []string `json:"roles"`
}

func toBusUpdateUser(u UpdateUser) (userbus.UpdateUser, error) {
	var errors errs.FieldErrors

	var firstName *name.Name

	if u.FirstName != nil {
		fName, err := name.Parse(*u.FirstName)
		if err != nil {
			errors.Add("first_name", err)
		}

		firstName = &fName
	}

	var lastName *name.Name

	if u.LastName != nil {
		lName, err := name.Parse(*u.LastName)
		if err != nil {
			errors.Add("last_name", err)
		}

		lastName = &lName
	}

	var pass *password.Password

	if u.Password != nil {
		ps, err := password.Parse(*u.Password)
		if err != nil {
			errors.Add("password", err)
		}

		pass = &ps
	}

	var email *mail.Address

	if u.Email != nil {
		addr, err := mail.ParseAddress(*u.Email)
		if err != nil {
			errors.Add("email", err)
		}

		email = addr
	}

	var roles []role.Role

	if u.Roles != nil {
		r, err := role.ParseMany(u.Roles)
		if err != nil {
			errors.Add("roles", err)
		}

		roles = r
	}

	if len(errors) > 0 {
		return userbus.UpdateUser{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	return userbus.UpdateUser{
		FirstName: firstName,
		LastName:  lastName,
		Password:  pass,
		Email:     email,
		Roles:     roles,
	}, nil
}
