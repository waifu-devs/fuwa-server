// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	ulid "github.com/oklog/ulid/v2"
)

type Channel struct {
	ChannelID ulid.ULID        `json:"channel_id"`
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type ChannelMessage struct {
	MessageID ulid.ULID        `json:"message_id"`
	ChannelID ulid.ULID        `json:"channel_id"`
	AuthorID  []byte           `json:"author_id"`
	Content   pgtype.Text      `json:"content"`
	Timestamp pgtype.Timestamp `json:"timestamp"`
}

type Config struct {
	ID     string      `json:"id"`
	Value  pgtype.Text `json:"value"`
	Public bool        `json:"public"`
}

type Role struct {
	RoleID ulid.ULID `json:"role_id"`
	Name   string    `json:"name"`
	Color  int32     `json:"color"`
}

type User struct {
	UserID     ulid.ULID        `json:"user_id"`
	FuwaUserID ulid.ULID        `json:"fuwa_user_id"`
	JoinedAt   pgtype.Timestamp `json:"joined_at"`
}

type UserRole struct {
	UserID ulid.ULID `json:"user_id"`
	RoleID ulid.ULID `json:"role_id"`
}

type UserSession struct {
	SessionID ulid.ULID        `json:"session_id"`
	UserID    ulid.ULID        `json:"user_id"`
	ExpiresAt pgtype.Timestamp `json:"expires_at"`
}
