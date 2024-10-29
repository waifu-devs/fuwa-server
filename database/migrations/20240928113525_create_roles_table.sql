-- +goose Up
CREATE TABLE roles (
	role_id BLOB NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	color INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE user_roles (
	user_id BLOB NOT NULL,
	role_id BLOB NOT NULL,

	CONSTRAINT user_roles_pk PRIMARY KEY (user_id, role_id),

	CONSTRAINT user_roles_user_id_fk
	FOREIGN KEY (user_id)
	REFERENCES users (user_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE,

	CONSTRAINT user_roles_role_id_fk
	FOREIGN KEY (role_id)
	REFERENCES roles (role_id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE user_roles;

DROP TABLE roles;
