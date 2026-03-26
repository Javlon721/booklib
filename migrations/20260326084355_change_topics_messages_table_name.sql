-- +goose Up
ALTER TABLE topics_messages
RENAME TO discussions;

-- +goose Down
ALTER TABLE discussions
RENAME TO topics_messages;