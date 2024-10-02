-- +goose Up
CREATE TABLE channel_messages (
	message_id BYTEA NOT NULL PRIMARY KEY,
	channel_id BYTEA NOT NULL,
	author_id BYTEA,
	content TEXT,
	timestamp TIMESTAMP NOT NULL
);

ALTER TABLE channel_messages
ADD CONSTRAINT channel_messages_channel_id_fk
FOREIGN KEY (channel_id)
REFERENCES channels (channel_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

-- +goose Down
ALTER TABLE channel_messages
DROP CONSTRAINT channel_messages_channel_id_fk;

DROP TABLE channel_messages;
