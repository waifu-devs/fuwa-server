-- +goose Up
CREATE TABLE channels (
	channel_id BYTEA NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE channels;
