package store

import (
	"context"
	"errors"
	"log/slog"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
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

func (s Store) CreateChat(ctx context.Context, nch chatbus.Chat) (uuid.UUID, error) {
	q := `
		INSERT INTO chats
			(user1, user2, date_updated, date_created)
		VALUES
			($1, $2, $3, $4)
		RETURNING chat_id;
	`

	var chatID uuid.UUID

	dbChat := toDBChat(nch)

	row := s.db.QueryRow(ctx, q, dbChat.User1, dbChat.User2, dbChat.DateUpdated, dbChat.DateCreated)

	if err := sqldb.ExtractPosgreErr(row.Scan(&chatID)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return uuid.UUID{}, chatbus.ErrChatAlreadyExists
		}
		return uuid.UUID{}, err
	}

	return chatID, nil
}

func (s Store) GetChatBy(ctx context.Context, user1, user2 uuid.UUID) (chatbus.Chat, error) {
	q := `
		SELECT chat_id, date_updated, date_created
			FROM chats
		WHERE $1 in (user1, user2) and $2 in (user1, user2);
	`

	chat := chatDB{
		User1: user1,
		User2: user2,
	}

	row := s.db.QueryRow(ctx, q, user1, user2)

	if err := sqldb.ExtractPosgreErr(row.Scan(&chat.ID, &chat.DateUpdated, &chat.DateCreated)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return chatbus.Chat{}, chatbus.ErrChatNotFound
		}
		return chatbus.Chat{}, err
	}

	return toBusChat(chat)
}

func (s Store) GetChatByID(ctx context.Context, chatID uuid.UUID) (chatbus.Chat, error) {
	q := `
		SELECT user1, user2, date_updated, date_created
			FROM chats
		WHERE chat_id = $1;
	`

	chat := chatDB{
		ID: chatID,
	}

	row := s.db.QueryRow(ctx, q, chatID)

	if err := sqldb.ExtractPosgreErr(row.Scan(&chat.User1, &chat.User2, &chat.DateUpdated, &chat.DateCreated)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return chatbus.Chat{}, chatbus.ErrChatNotFound
		}
		return chatbus.Chat{}, err
	}

	return toBusChat(chat)
}

func (s Store) CreateMessage(ctx context.Context, nm chatbus.Message) (uuid.UUID, error) {
	q := `
		INSERT INTO messages
			(chat_id, sender, content, status, date_updated, date_created)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	var id uuid.UUID

	dbMessage := toDBmessage(nm)

	row := s.db.QueryRow(ctx, q, dbMessage.ChatID, dbMessage.Sender, dbMessage.Content, dbMessage.Status, dbMessage.DateCreated, dbMessage.DateUpdated)

	if err := sqldb.ExtractPosgreErr(row.Scan(&id)); err != nil {
		if errors.Is(err, sqldb.ErrFfreignKeyViolation) {
			return uuid.UUID{}, chatbus.ErrChatNotFound
		}
		return uuid.UUID{}, err
	}

	return id, nil
}

func (s Store) GetMessagesByStatus(ctx context.Context, chatID, reciever uuid.UUID, status string) ([]chatbus.Message, error) {
	q := `
		SELECT id, chat_id, content, sender, date_created, date_updated, status
			FROM messages
		WHERE chat_id = $1 and status = $2 and sender <> $3;
	`

	rows, err := s.db.Query(ctx, q, chatID, status, reciever)

	if err != nil {
		return []chatbus.Message{}, err
	}

	defer rows.Close()

	messages, err := pgx.CollectRows(rows, pgx.RowToStructByName[messageDB])

	if err != nil {
		return []chatbus.Message{}, err
	}

	return toBusMessages(messages)
}
