-- +goose Up
CREATE TABLE config (
	config_id TEXT NOT NULL PRIMARY KEY,
	value TEXT,
	public INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE config;
