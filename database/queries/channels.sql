-- name: ListChannelsAll :many
SELECT * FROM channels;

-- name: ListChannels :many
SELECT * FROM channels
WHERE channel_id > $1 LIMIT $2;

-- name: GetChannel :one
SELECT * FROM channels
WHERE channel_id = $1;

-- name: CreateChannel :exec
INSERT INTO channels (channel_id, name, type, created_at)
VALUES ($1, $2, $3, $4);
