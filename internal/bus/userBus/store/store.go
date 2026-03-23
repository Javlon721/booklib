package store

import (
	"context"
	"errors"
	"log/slog"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	sqldb "github.com/Javlon721/booklib/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	db     sqldb.DB
	logger *slog.Logger
}

func New(db sqldb.DB, logger *slog.Logger) *Store {
	return &Store{
		db:     db,
		logger: logger,
	}
}

func (s Store) GetUserByID(ctx context.Context, userID uuid.UUID) (userbus.User, error) {
	const q = `
	SELECT user_id, first_name, last_name, email, hash_password, date_created, date_updated, roles
	FROM users 
	WHERE user_id = $1
	`

	var dbUser userDB

	row := s.db.QueryRow(ctx, q, userID)

	if err := sqldb.ExtractPosgreErr(row.Scan(&dbUser.ID, &dbUser.FirstName, &dbUser.LastName, &dbUser.Email, &dbUser.PasswordHash, &dbUser.DateCreated, &dbUser.DateUpdated, &dbUser.Roles)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userbus.User{}, userbus.ErrUserNotFound
		}
		return userbus.User{}, err
	}

	return toBusUser(dbUser)
}

func (s Store) GetUserByEmail(ctx context.Context, email string) (userbus.User, error) {
	const q = `
	SELECT user_id, first_name, last_name, email, hash_password, date_created, date_updated, roles
	FROM users 
	WHERE email = $1
	`

	var dbUser userDB

	row := s.db.QueryRow(ctx, q, email)

	if err := sqldb.ExtractPosgreErr(row.Scan(&dbUser.ID, &dbUser.FirstName, &dbUser.LastName, &dbUser.Email, &dbUser.PasswordHash, &dbUser.DateCreated, &dbUser.DateUpdated, &dbUser.Roles)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userbus.User{}, userbus.ErrUserNotFound
		}
		return userbus.User{}, err
	}

	return toBusUser(dbUser)
}

func (s Store) Create(ctx context.Context, nu userbus.User) (uuid.UUID, error) {
	const q = `
	INSERT INTO users 
		(first_name, last_name, email, hash_password, date_created, date_updated, roles)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING user_id`

	var id uuid.UUID

	dbUser := toDBUser(nu)

	s.logger.Info("create user", "roles", dbUser.Roles)

	row := s.db.QueryRow(ctx, q, dbUser.FirstName, dbUser.LastName, dbUser.Email, dbUser.PasswordHash, dbUser.DateCreated, dbUser.DateUpdated, dbUser.Roles)

	if err := sqldb.ExtractPosgreErr(row.Scan(&id)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return id, userbus.ErrUserAlreadyExists
		}
		return id, err
	}

	return id, nil
}
