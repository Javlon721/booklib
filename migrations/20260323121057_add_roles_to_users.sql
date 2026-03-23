-- +goose Up
ALTER TABLE users
ADD COLUMN roles varchar[] not null;

-- +goose Down
ALTER TABLE users
DROP COLUMN roles;
