-- +goose Up
CREATE TABLE refresh_tokens (
  token VARCHAR PRIMARY KEY,
  user_id UUID REFERENCES users ON DELETE CASCADE NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE refresh_tokens;
