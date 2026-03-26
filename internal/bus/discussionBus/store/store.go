package store

import (
	"context"
	"log/slog"

	discussionBus "github.com/Javlon721/booklib/internal/bus/discussionBus"
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

func (s Store) Create(ctx context.Context, ntm discussionBus.TopicMessage) (uuid.UUID, error) {
	q := `
		INSERT INTO discussions
			(message, date_created, date_updated, user_id, topic_id)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id;
	`

	dbTopicMessage := toDBTopicMessage(ntm)

	var id uuid.UUID

	row := s.db.QueryRow(ctx, q, dbTopicMessage.Message, dbTopicMessage.DateCreated, dbTopicMessage.DateUpdated, dbTopicMessage.UserID, dbTopicMessage.TopicID)

	if err := sqldb.ExtractPosgreErr(row.Scan(&id)); err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
