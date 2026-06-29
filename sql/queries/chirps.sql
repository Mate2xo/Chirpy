-- name: CreateChirp :one
INSERT INTO chirps (id, body, user_id, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  NOW(),
  NOW()
)
RETURNING *;

-- name: AllChirps :many
SELECT * FROM chirps
  WHERE (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id')::uuid)
  ORDER BY
    CASE WHEN NOT @reverse::boolean THEN created_at END ASC,
    CASE WHEN @reverse::boolean  THEN created_at END  DESC;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChrip :exec
DELETE FROM chirps WHERE id = $1;
