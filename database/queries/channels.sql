-- name: ListChannelsAll :many
SELECT * FROM channels;

-- name: ListChannels :many
SELECT * FROM channels
WHERE channel_id > ? LIMIT ?;

-- name: GetChannel :one
SELECT * FROM channels
WHERE channel_id = ?;

-- name: CreateChannel :exec
INSERT INTO channels (channel_id, name, type, created_at)
VALUES (?, ?, ?, ?);

-- name: UpdateChannel :exec
UPDATE channels
SET name = ?, type = ?
WHERE channel_id = ?;
