package model

type User struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Role     int32  `db:"role"`
}
