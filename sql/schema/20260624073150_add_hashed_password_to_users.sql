-- +goose Up
SET LOCAL lock_timeout = '2s';
ALTER TABLE users
  ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
ALTER TABLE users DROP COLUMN hashed_password;
