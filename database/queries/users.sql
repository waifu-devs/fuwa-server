-- name: ListUsers :many
SELECT * FROM users
WHERE user_id > ? LIMIT ?;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = ?;

-- name: CreateUser :exec
INSERT INTO users (user_id, fuwa_user_id, joined_at) VALUES (?, ?, ?);
