package userbus

import (
	"net/mail"
	"time"

	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/Javlon721/booklib/internal/types/password"
	"github.com/Javlon721/booklib/internal/types/role"
	"github.com/google/uuid"
)

type CreateUser struct {
	FirstName name.Name
	LastName  name.Name
	Password  password.Password
	Email     mail.Address
	Roles     []role.Role
}

type User struct {
	ID           uuid.UUID
	FirstName    name.Name
	LastName     name.Name
	Email        mail.Address
	PasswordHash []byte
	DateCreated  time.Time
	DateUpdated  time.Time
	Roles        []role.Role
}
