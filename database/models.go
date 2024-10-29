// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"

	ulid "github.com/oklog/ulid/v2"
)

type Channel struct {
	ChannelID ulid.ULID `json:"channel_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt int64     `json:"created_at"`
}

type ChannelMessage struct {
	MessageID ulid.ULID      `json:"message_id"`
	ChannelID ulid.ULID      `json:"channel_id"`
	AuthorID  ulid.ULID      `json:"author_id"`
	Content   sql.NullString `json:"content"`
	Timestamp int64          `json:"timestamp"`
}

type Config struct {
	ConfigID string         `json:"config_id"`
	Value    sql.NullString `json:"value"`
	Public   int64          `json:"public"`
}

type Role struct {
	RoleID ulid.ULID `json:"role_id"`
	Name   string    `json:"name"`
	Color  int64     `json:"color"`
}

type User struct {
	UserID     ulid.ULID `json:"user_id"`
	FuwaUserID ulid.ULID `json:"fuwa_user_id"`
	JoinedAt   int64     `json:"joined_at"`
}

type UserRole struct {
	UserID ulid.ULID `json:"user_id"`
	RoleID ulid.ULID `json:"role_id"`
}

type UserSession struct {
	SessionID ulid.ULID `json:"session_id"`
	UserID    ulid.ULID `json:"user_id"`
	ExpiresAt int64     `json:"expires_at"`
}
