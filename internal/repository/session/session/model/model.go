package model

import "time"

type Session struct {
	ID           int64      `db:"id"`
	UserID       int64      `db:"user_id"`
	RefreshToken string     `db:"refresh_token"`
	CreatedAt    time.Time  `db:"created_at"`
	ExpiresAt    time.Time  `db:"expires_at"`
	RevokedAt    *time.Time `db:"revoked_at"`
}
