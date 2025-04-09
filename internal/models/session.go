package models

import "time"

type Session struct {
	ID           int64
	UserID       int64
	RefreshToken string
	CreatedAt    time.Time
	ExpiresAt    time.Time
	RevokedAt    *time.Time
}
