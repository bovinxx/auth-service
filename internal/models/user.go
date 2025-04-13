package models

type Role int

const (
	RoleAdmin = iota
	RoleUser
)

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Role     Role
}
