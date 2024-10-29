-- +goose Up
CREATE TABLE users (
	user_id BLOB NOT NULL PRIMARY KEY,
	fuwa_user_id BLOB NOT NULL,
	joined_at INTEGER NOT NULL
);

-- +goose Down
DROP TABLE users;
