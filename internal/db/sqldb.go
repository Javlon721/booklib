package sqldb

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolation = "23505"
	undefinedTable  = "42P01"
)

var (
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
	ErrUndefinedTable    = errors.New("undefined table")
)

type DB interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

func ExtractPosgreErr(err error) error {
	if err == nil {
		return nil
	}

	var pqerr *pgconn.PgError

	if errors.As(err, &pqerr) {
		switch pqerr.Code {
		case undefinedTable:
			return ErrUndefinedTable
		case uniqueViolation:
			return ErrDBDuplicatedEntry
		}
	}

	return err
}
