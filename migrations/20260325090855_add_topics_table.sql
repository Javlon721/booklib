-- +goose Up
CREATE TABLE
  IF NOT EXISTS topics (
    topic_id uuid DEFAULT gen_random_uuid (),
    title VARCHAR NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,
    created_by uuid NOT NULL,
    PRIMARY KEY (topic_id),
    FOREIGN KEY (created_by) REFERENCES users (user_id)
  );

CREATE TABLE
  IF NOT EXISTS topics_messages (
    id uuid DEFAULT gen_random_uuid (),
    message VARCHAR NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,
    user_id uuid NOT NULL,
    topic_id uuid NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (user_id),
    FOREIGN KEY (topic_id) REFERENCES topics (topic_id)
  );

-- +goose Down
DROP TABLE IF EXISTS topics;

DROP TABLE IF EXISTS topics_messages;