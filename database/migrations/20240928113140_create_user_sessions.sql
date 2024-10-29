-- +goose Up
CREATE TABLE user_sessions (
	session_id BLOB NOT NULL PRIMARY KEY,
	user_id BLOB NOT NULL,
	expires_at INTEGER NOT NULL,

	CONSTRAINT user_sessions_user_id_fk
	FOREIGN KEY (user_id)
	REFERENCES users (user_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE user_sessions;
