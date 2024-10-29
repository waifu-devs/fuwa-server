// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: servers.sql

package database

import (
	"context"

	ulid "github.com/oklog/ulid/v2"
)

const createServer = `-- name: CreateServer :exec
INSERT INTO servers (server_id, name)
VALUES (?, ?)
`

type CreateServerParams struct {
	ServerID ulid.ULID `json:"server_id"`
	Name     string    `json:"name"`
}

func (q *Queries) CreateServer(ctx context.Context, arg CreateServerParams) error {
	_, err := q.db.ExecContext(ctx, createServer, arg.ServerID, arg.Name)
	return err
}

const createServerUser = `-- name: CreateServerUser :exec
INSERT INTO server_users (server_user_id, server_id, user_id, display_name)
VALUES (?, ?, ?, ?)
`

type CreateServerUserParams struct {
	ServerUserID ulid.ULID `json:"server_user_id"`
	ServerID     ulid.ULID `json:"server_id"`
	UserID       ulid.ULID `json:"user_id"`
	DisplayName  string    `json:"display_name"`
}

func (q *Queries) CreateServerUser(ctx context.Context, arg CreateServerUserParams) error {
	_, err := q.db.ExecContext(ctx, createServerUser,
		arg.ServerUserID,
		arg.ServerID,
		arg.UserID,
		arg.DisplayName,
	)
	return err
}

const deleteServerUser = `-- name: DeleteServerUser :exec
DELETE FROM server_users
WHERE server_user_id = ?
`

func (q *Queries) DeleteServerUser(ctx context.Context, serverUserID ulid.ULID) error {
	_, err := q.db.ExecContext(ctx, deleteServerUser, serverUserID)
	return err
}

const listServerUsers = `-- name: ListServerUsers :many
SELECT server_user_id, server_id, user_id, display_name FROM server_users
WHERE server_user_id > ? LIMIT ?
`

type ListServerUsersParams struct {
	ServerUserID ulid.ULID `json:"server_user_id"`
	Limit        int64     `json:"limit"`
}

func (q *Queries) ListServerUsers(ctx context.Context, arg ListServerUsersParams) ([]ServerUser, error) {
	rows, err := q.db.QueryContext(ctx, listServerUsers, arg.ServerUserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ServerUser
	for rows.Next() {
		var i ServerUser
		if err := rows.Scan(
			&i.ServerUserID,
			&i.ServerID,
			&i.UserID,
			&i.DisplayName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listServerUsersAll = `-- name: ListServerUsersAll :many
SELECT server_user_id, server_id, user_id, display_name FROM server_users
`

func (q *Queries) ListServerUsersAll(ctx context.Context) ([]ServerUser, error) {
	rows, err := q.db.QueryContext(ctx, listServerUsersAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ServerUser
	for rows.Next() {
		var i ServerUser
		if err := rows.Scan(
			&i.ServerUserID,
			&i.ServerID,
			&i.UserID,
			&i.DisplayName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listServers = `-- name: ListServers :many
SELECT server_id, name FROM servers
WHERE server_id > ? LIMIT ?
`

type ListServersParams struct {
	ServerID ulid.ULID `json:"server_id"`
	Limit    int64     `json:"limit"`
}

func (q *Queries) ListServers(ctx context.Context, arg ListServersParams) ([]Server, error) {
	rows, err := q.db.QueryContext(ctx, listServers, arg.ServerID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Server
	for rows.Next() {
		var i Server
		if err := rows.Scan(&i.ServerID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listServersAll = `-- name: ListServersAll :many
SELECT server_id, name FROM servers
`

func (q *Queries) ListServersAll(ctx context.Context) ([]Server, error) {
	rows, err := q.db.QueryContext(ctx, listServersAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Server
	for rows.Next() {
		var i Server
		if err := rows.Scan(&i.ServerID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateServer = `-- name: UpdateServer :exec
UPDATE servers
SET name = ?
WHERE server_id = ?
`

type UpdateServerParams struct {
	Name     string    `json:"name"`
	ServerID ulid.ULID `json:"server_id"`
}

func (q *Queries) UpdateServer(ctx context.Context, arg UpdateServerParams) error {
	_, err := q.db.ExecContext(ctx, updateServer, arg.Name, arg.ServerID)
	return err
}

const updateServerUser = `-- name: UpdateServerUser :exec
UPDATE server_users
SET display_name = ?
WHERE server_user_id = ?
`

type UpdateServerUserParams struct {
	DisplayName  string    `json:"display_name"`
	ServerUserID ulid.ULID `json:"server_user_id"`
}

func (q *Queries) UpdateServerUser(ctx context.Context, arg UpdateServerUserParams) error {
	_, err := q.db.ExecContext(ctx, updateServerUser, arg.DisplayName, arg.ServerUserID)
	return err
}