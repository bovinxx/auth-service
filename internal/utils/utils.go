package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(hashedPassword, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}

func NewCacheKey(prefix, username string) string {
	return fmt.Sprintf("%s:%s", prefix, username)
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
