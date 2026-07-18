package auth

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken() string {
	slice := make([]byte, 128)
	rand.Read(slice)
	return hex.EncodeToString(slice)
}

func CheckPassword (input string, hash string, logger *zerolog.Logger) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	if err != nil {
		logger.Error().Err(err).Msg("Error checking password")
		return false
	}
	return true
}
