-- +goose NO TRANSACTION

-- +goose Up
CREATE INDEX user_roles_server_user_id_idx
ON user_roles (user_id);

CREATE INDEX channel_messages_channel_id_idx
ON channel_messages (channel_id);
-- +goose Down
DROP INDEX user_roles_user_id_idx;

DROP INDEX channel_messages_channel_id_idx;
