-- +goose Up
CREATE TABLE channel_messages (
	message_id BLOB NOT NULL CONSTRAINT channel_messages_pk PRIMARY KEY,
	server_id BLOB NOT NULL,
	channel_id BLOB NOT NULL,
	author_id BLOB NOT NULL,
	content TEXT,
	timestamp INTEGER NOT NULL,

	CONSTRAINT channel_messages_server_id_fk
	FOREIGN KEY (server_id)
	REFERENCES servers (server_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE,

	CONSTRAINT channel_messages_channel_id_fk
	FOREIGN KEY (server_id, channel_id)
	REFERENCES channels (server_id, channel_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE,

	CONSTRAINT channel_messages_author_id_fk
	FOREIGN KEY (author_id)
	REFERENCES server_users (server_user_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE channel_messages;
