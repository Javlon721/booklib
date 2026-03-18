-- +goose Up
CREATE TABLE
  IF NOT EXISTS users (
    user_id uuid DEFAULT gen_random_uuid(),
    email VARCHAR(150) UNIQUE NOT NULL,
    hash_password VARCHAR NOT NULL,
    first_name VARCHAR(150) NOT NULL,
    last_name VARCHAR(150),
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,

    PRIMARY KEY (user_id)
  );

-- +goose Down
DROP TABLE IF EXISTS users;