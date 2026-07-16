package auth

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func GenerateToken() string {
	slice := make([]byte, 128)
	rand.Read(slice)
	return hex.EncodeToString(slice)
}

func CheckPassword (input string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(input), []byte(hash))
	if err != nil {
		return false
	}
	return true
}
