-- +goose Up
CREATE TABLE channels (
	channel_id BLOB NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	created_at INTEGER NOT NULL
);


-- +goose Down
DROP TABLE channels;
