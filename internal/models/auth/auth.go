package auth

import "github.com/golang-jwt/jwt/v4"

type UserInfo struct {
	Login string
	Role  string
}

type UserClaims struct {
	jwt.StandardClaims
	Username string
	Role     string
}
