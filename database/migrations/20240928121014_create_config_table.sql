-- +goose Up
CREATE TABLE config (
	id TEXT NOT NULL PRIMARY KEY,
	value TEXT,
	public BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE config;
