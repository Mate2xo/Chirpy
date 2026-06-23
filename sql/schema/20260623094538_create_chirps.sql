-- +goose Up
CREATE TABLE chirps (
  id UUID PRIMARY KEY,
  body TEXT NOT NULL,
  user_id UUID REFERENCES users ON DELETE CASCADE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE chirps;
