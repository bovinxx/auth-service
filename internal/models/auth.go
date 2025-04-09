package models

import "github.com/golang-jwt/jwt/v4"

type UserInfo struct {
	UserID   int64
	Username string
	Role     string
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID   int64
	Username string
	Role     string
}
