package jwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (j *JWT) Generate(payload *JWTPayload) (string, int64, error) {
	expiresIn, err := strconv.Atoi(payload.ExpiresAt)
	if err != nil {
		return "", 0, err
	}

	secretKey := []byte(j.config.JWTSecret)
	ExpiresToken := time.Now().Add(time.Hour * 24 * time.Duration(expiresIn)).Unix()

	claims := jwt.MapClaims{
		"exp": ExpiresToken,
		"iss": j.config.Name,
		"uid": payload.UID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, err
	}

	return t, ExpiresToken, nil
}
