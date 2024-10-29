-- name: ListChannelsAll :many
SELECT * FROM channels WHERE server_id = ?;

-- name: ListChannels :many
SELECT * FROM channels
WHERE server_id = ? AND channel_id > ? LIMIT ?;

-- name: GetChannel :one
SELECT * FROM channels
WHERE server_id = ? AND channel_id = ?;

-- name: CreateChannel :exec
INSERT INTO channels (server_id, channel_id, name, type, created_at)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateChannel :exec
UPDATE channels
SET name = ?, type = ?
WHERE server_id = ? AND channel_id = ?;
