-- +goose Up
CREATE TABLE users (
	user_id BYTEA NOT NULL PRIMARY KEY,
	fuwa_user_id BYTEA NOT NULL,
	joined_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;
