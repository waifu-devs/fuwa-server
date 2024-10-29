-- +goose Up
CREATE TABLE roles (
	server_id BLOB NOT NULL,
	role_id BLOB NOT NULL,
	name TEXT NOT NULL,
	color INTEGER NOT NULL DEFAULT 0,

	CONSTRAINT roles_pk PRIMARY KEY (server_id, role_id)
);

CREATE TABLE user_roles (
	server_id BLOB NOT NULL,
	server_user_id BLOB NOT NULL,
	role_id BLOB NOT NULL,

	CONSTRAINT user_roles_pk PRIMARY KEY (server_user_id, role_id),

	CONSTRAINT user_roles_server_user_id_fk
	FOREIGN KEY (server_user_id)
	REFERENCES server_users (server_user_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE,

	CONSTRAINT user_roles_role_id_fk
	FOREIGN KEY (server_id, role_id)
	REFERENCES roles (server_id, role_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE user_roles;

DROP TABLE roles;
