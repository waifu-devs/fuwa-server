-- +goose NO TRANSACTION

-- +goose Up
CREATE INDEX roles_server_id_idx
ON roles (server_id);

CREATE INDEX user_roles_server_id_idx
ON user_roles (server_id);

CREATE INDEX user_roles_server_user_id_idx
ON user_roles (server_user_id);

CREATE INDEX server_users_server_id_idx
ON server_users (server_id);

CREATE INDEX channels_server_id_idx
ON channels (server_id);

CREATE INDEX channel_messages_channel_id_idx
ON channel_messages (channel_id);
-- +goose Down
DROP INDEX roles_server_id_idx;

DROP INDEX user_roles_server_id_idx;

DROP INDEX user_roles_user_id_idx;

DROP INDEX server_users_server_id_idx;

DROP INDEX channels_server_id_idx;

DROP INDEX channel_messages_channel_id_idx;
