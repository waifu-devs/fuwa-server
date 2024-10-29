-- name: ListServersAll :many
SELECT * FROM servers;

-- name: ListServers :many
SELECT * FROM servers
WHERE server_id > ? LIMIT ?;

-- name: ListServerUsersAll :many
SELECT * FROM server_users;

-- name: ListServerUsers :many
SELECT * FROM server_users
WHERE server_user_id > ? LIMIT ?;

-- name: CreateServer :exec
INSERT INTO servers (server_id, name)
VALUES (?, ?);

-- name: UpdateServer :exec
UPDATE servers
SET name = ?
WHERE server_id = ?;

-- name: CreateServerUser :exec
INSERT INTO server_users (server_user_id, server_id, user_id, display_name)
VALUES (?, ?, ?, ?);

-- name: UpdateServerUser :exec
UPDATE server_users
SET display_name = ?
WHERE server_user_id = ?;

-- name: DeleteServerUser :exec
DELETE FROM server_users
WHERE server_user_id = ?;
