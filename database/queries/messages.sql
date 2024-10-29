-- name: ListMessages :many
SELECT * FROM channel_messages
WHERE server_id = ? AND channel_id = ? AND message_id > ? LIMIT ?;

-- name: GetMessage :one
SELECT * FROM channel_messages
WHERE message_id = ?;

-- name: CreateMessage :exec
INSERT INTO channel_messages (
	message_id,
	server_id,
	channel_id,
	author_id,
	content,
	timestamp
)
VALUES (?, ?, ?, ?, ?, ?);
