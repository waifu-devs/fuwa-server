-- name: ListUsers :many
SELECT * FROM users
WHERE user_id > $1 LIMIT $2;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1;

-- name: CreateUser :exec
INSERT INTO users (user_id, fuwa_user_id, joined_at) VALUES ($1, $2, $3);
