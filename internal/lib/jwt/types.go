package jwt

import "gofi/internal/config"

type JWT struct {
	config *config.ConfigApp
}

type JWTPayload struct {
	UID       string
	Secret    string
	ExpiresAt string
}

type JWTClaims struct {
	Exp int64
	Iss string
	UID string
}
