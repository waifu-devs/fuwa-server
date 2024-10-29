-- +goose Up
CREATE TABLE servers (
	server_id BLOB NOT NULL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE server_users (
	server_user_id BLOB NOT NULL PRIMARY KEY,
	server_id BLOB NOT NULL,
	user_id BLOB NOT NULL,
	display_name TEXT NOT NULL,

	CONSTRAINT server_users_server_id_fk
	FOREIGN KEY (server_id)
	REFERENCES servers (server_id)
	ON DELETE CASCADE ON UPDATE CASCADE,

	CONSTRAINT server_users_user_id_fk
	FOREIGN KEY (user_id)
	REFERENCES users (user_id)
	ON DELETE CASCADE ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE server_users;

DROP TABLE servers;
