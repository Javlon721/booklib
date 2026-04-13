package userbus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Store interface {
	Create(context.Context, User) (uuid.UUID, error)
	GetUserByID(context.Context, uuid.UUID) (User, error)
	GetUserByEmail(context.Context, string) (User, error)
	Update(context.Context, User) error
}

type Business struct {
	logger *slog.Logger
	store  Store
}

func NewBusiness(logger *slog.Logger, store Store) *Business {
	return &Business{
		store:  store,
		logger: logger,
	}
}

func (bus Business) Create(ctx context.Context, nu CreateUser) (User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(nu.Password.String()), bcrypt.DefaultCost)

	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	user := User{
		FirstName:    nu.FirstName,
		LastName:     nu.LastName,
		PasswordHash: passwordHash,
		Email:        nu.Email,
		DateCreated:  now,
		DateUpdated:  now,
		Roles:        nu.Roles,
	}

	id, err := bus.store.Create(ctx, user)

	if err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	user.ID = id

	return user, nil
}

func (bus Business) GetUserByID(ctx context.Context, userID uuid.UUID) (User, error) {
	return bus.store.GetUserByID(ctx, userID)
}

func (bus Business) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return bus.store.GetUserByEmail(ctx, email)
}
func (bus Business) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
	bus.logger.Info("bus", "updateUser", uu)

	if uu.FirstName != nil {
		usr.FirstName = *uu.FirstName
	}

	if uu.LastName != nil {
		usr.LastName = *uu.LastName
	}

	if uu.Email != nil {
		usr.Email = *uu.Email
	}

	if uu.Password != nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(uu.Password.String()), bcrypt.DefaultCost)

		if err != nil {
			return User{}, fmt.Errorf("generatefrompassword: %w", err)
		}

		usr.PasswordHash = passwordHash
	}

	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}

	usr.DateUpdated = time.Now()

	err := bus.store.Update(ctx, usr)

	if err != nil {
		return User{}, err
	}

	return usr, nil
}
