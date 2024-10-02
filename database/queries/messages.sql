-- name: ListMessages :many
SELECT * FROM channel_messages
WHERE channel_id = $1 AND message_id > $2 LIMIT $3;

-- name: GetMessage :one
SELECT * FROM channel_messages
WHERE message_id = $1;

-- name: CreateMessage :exec
INSERT INTO channel_messages (
	message_id,
	channel_id,
	author_id,
	content,
	timestamp
)
VALUES ($1, $2, $3, $4, $5);
