package store

import (
	"context"
	"log/slog"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	sqldb "github.com/Javlon721/booklib/internal/db"
	"github.com/google/uuid"
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
	SELECT user_id, first_name, last_name, email, hash_password, date_created, date_updated
	FROM users 
	WHERE user_id = $1
	`

	var dbUser userDB

	row := s.db.QueryRow(ctx, q, userID)

	if err := sqldb.ExtractPosgreErr(row.Scan(&dbUser.ID, &dbUser.FirstName, &dbUser.LastName, &dbUser.Email, &dbUser.PasswordHash, &dbUser.DateCreated, &dbUser.DateUpdated)); err != nil {
		return userbus.User{}, err
	}

	return toBusUser(dbUser)
}

func (s Store) Create(ctx context.Context, nu userbus.User) (uuid.UUID, error) {
	const q = `
	INSERT INTO users 
		(first_name, last_name, email, hash_password, date_created, date_updated)
	VALUES 
		($1, $2, $3, $4, $5, $6)
	RETURNING user_id`

	var id uuid.UUID

	dbUser := toDBUser(nu)

	row := s.db.QueryRow(ctx, q, dbUser.FirstName, dbUser.LastName, dbUser.Email, dbUser.PasswordHash, dbUser.DateCreated, dbUser.DateUpdated)

	if err := sqldb.ExtractPosgreErr(row.Scan(&id)); err != nil {
		return id, err
	}

	return id, nil
}
