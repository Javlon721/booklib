-- TODO i know how shitty this is, so later i will improve it
-- +goose Up
CREATE TABLE
  IF NOT EXISTS chats (
    chat_id uuid DEFAULT gen_random_uuid (),
    user1 uuid NOT NULL,
    user2 uuid NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,
    PRIMARY KEY (chat_id),
    FOREIGN KEY (user1) REFERENCES users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (user2) REFERENCES users (user_id) ON DELETE CASCADE,
    UNIQUE (user1, user2)
  );

CREATE TABLE
  IF NOT EXISTS messages (
    id uuid DEFAULT gen_random_uuid (),
    chat_id uuid NOT NULL,
    sender uuid NOT NULL,
    content TEXT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,
    status VARCHAR NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (chat_id) REFERENCES chats (chat_id) ON DELETE CASCADE,
    FOREIGN KEY (sender) REFERENCES users (user_id) ON DELETE CASCADE
  );

-- +goose Down
DROP TABLE IF EXISTS chats;

DROP TABLE IF EXISTS messages;