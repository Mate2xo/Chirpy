-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
RETURNING id, email, hashed_password, created_at;

-- name: UpdateUser :one
UPDATE users
  SET hashed_password = $1, email = $2, updated_at = NOW()
  WHERE id = $3
RETURNING *;

-- name: UpgradeUser :exec
UPDATE users
  SET is_chirpy_red = true, updated_at = NOW()
  WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;
