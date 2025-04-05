package auth

import "github.com/golang-jwt/jwt/v4"

type UserInfo struct {
	Login string
	Role  string
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username string
	Role     string
}
