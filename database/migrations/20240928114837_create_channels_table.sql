-- +goose Up
CREATE TABLE channels (
	server_id BLOB NOT NULL,
	channel_id BLOB NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	created_at INTEGER NOT NULL,

	CONSTRAINT channels_pk PRIMARY KEY (server_id, channel_id),

	CONSTRAINT channels_server_id_fk
	FOREIGN KEY (server_id)
	REFERENCES servers (server_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);


-- +goose Down
DROP TABLE channels;
