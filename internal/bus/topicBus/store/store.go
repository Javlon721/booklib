package store

import (
	"bytes"
	"context"
	"errors"
	"log/slog"

	topicbus "github.com/Javlon721/booklib/internal/bus/topicBus"
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

func (s Store) Create(ctx context.Context, nt topicbus.Topic) (uuid.UUID, error) {
	q := `
		INSERT INTO topics
			(title, date_created, date_updated, created_by)
		VALUES
			($1, $2, $3, $4)
		RETURNING topic_id;
	`

	var topicID uuid.UUID
	row := s.db.QueryRow(ctx, q, nt.Title, nt.DateCreated, nt.DateUpdated, nt.CreatedBy)

	if err := sqldb.ExtractPosgreErr(row.Scan(&topicID)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return uuid.UUID{}, topicbus.ErrTopicAlreadyExists
		}
		return uuid.UUID{}, err
	}

	return topicID, nil
}

func (s Store) GetByID(ctx context.Context, topicID uuid.UUID) (topicbus.Topic, error) {
	q := `
		SELECT topic_id, title, date_created, date_updated, created_by
		FROM topics
		WHERE topic_id = $1;
	`

	row := s.db.QueryRow(ctx, q, topicID)

	var dbTopic topicDB

	if err := sqldb.ExtractPosgreErr(row.Scan(&dbTopic.ID, &dbTopic.Title, &dbTopic.DateCreated, &dbTopic.DateUpdated, &dbTopic.CreatedBy)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return topicbus.Topic{}, topicbus.ErrTopicNotFound
		}
		return topicbus.Topic{}, err
	}

	return toBusTopic(dbTopic)
}

func (s Store) Query(ctx context.Context, filter topicbus.QueryFilter) ([]topicbus.Topic, error) {
	q := `
		SELECT topic_id, title, date_created, date_updated, created_by
		FROM topics
	`

	buf := bytes.NewBufferString(q)
	namedArgs := applyFilter(filter, buf)

	s.logger.Info("store.query", "q", buf.String())

	rows, err := s.db.Query(ctx, buf.String(), namedArgs)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	topics, err := pgx.CollectRows(rows, pgx.RowToStructByName[topicDB])

	if err != nil {
		return nil, err
	}

	return toBusTopics(topics)
}
