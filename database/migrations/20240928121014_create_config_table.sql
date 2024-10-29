-- +goose Up
CREATE TABLE config (
	server_id BLOB NOT NULL,
	config_id TEXT NOT NULL,
	value TEXT,
	public INTEGER NOT NULL DEFAULT 0,

	CONSTRAINT config_pk PRIMARY KEY (server_id, config_id),

	CONSTRAINT config_server_id_fk
	FOREIGN KEY (server_id)
	REFERENCES servers (server_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE config;
