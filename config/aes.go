package config

import (
	"github.com/mervick/aes-everywhere/go/aes256"
)

func GenerateToken(data string, config Config) string {
	jwtSecret := config.JwtSecret

	tokenSigned := aes256.Encrypt(data, jwtSecret)
	return tokenSigned
}
