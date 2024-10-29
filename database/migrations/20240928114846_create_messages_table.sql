-- +goose Up
CREATE TABLE channel_messages (
	message_id BLOB NOT NULL CONSTRAINT channel_messages_pk PRIMARY KEY,
	channel_id BLOB NOT NULL,
	author_id BLOB,
	content TEXT NOT NULL,
	timestamp INTEGER NOT NULL,

	CONSTRAINT channel_messages_channel_id_fk
	FOREIGN KEY (channel_id)
	REFERENCES channels (channel_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE

	--CONSTRAINT channel_messages_author_id_fk
	--FOREIGN KEY (author_id)
	--REFERENCES server_users (user_id)
	--ON DELETE CASCADE
	--ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE channel_messages;
