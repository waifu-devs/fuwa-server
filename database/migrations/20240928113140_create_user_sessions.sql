-- +goose Up
CREATE TABLE user_sessions (
	session_id BYTEA NOT NULL PRIMARY KEY,
	user_id BYTEA NOT NULL,
	expires_at TIMESTAMP NOT NULL
);

ALTER TABLE user_sessions
ADD CONSTRAINT user_sessions_user_id_fk
FOREIGN KEY (user_id)
REFERENCES users (user_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

-- +goose Down
ALTER TABLE user_sessions
DROP CONSTRAINT user_sessions_user_id_fk;

DROP TABLE user_sessions;
