-- +goose Up
CREATE TABLE roles (
	role_id BYTEA NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	color INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE user_roles (
	user_id BYTEA NOT NULL,
	role_id BYTEA NOT NULL,
	PRIMARY KEY (user_id, role_id)
);

ALTER TABLE user_roles
ADD CONSTRAINT user_roles_user_id_fk
FOREIGN KEY (user_id)
REFERENCES users (user_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE user_roles
ADD CONSTRAINT user_roles_role_id_fk
FOREIGN KEY (role_id)
REFERENCES roles (role_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

-- +goose Down
ALTER TABLE user_roles
DROP CONSTRAINT user_roles_user_id_fk;

ALTER TABLE user_roles
DROP CONSTRAINT user_roles_role_id_fk;

DROP TABLE user_roles;

DROP TABLE roles;
