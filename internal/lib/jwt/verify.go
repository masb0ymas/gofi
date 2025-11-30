package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func (j *JWT) Verify(extractToken string) (*JWTClaims, error) {
	token, err := jwt.Parse(extractToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return []byte(j.config.JWTSecret), nil
	})

	if err != nil {
		message := fmt.Sprintf("Token verification failed: %v. Please ensure your token is valid and not expired.", err)
		return nil, errors.New(message)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &JWTClaims{
			Exp: int64(claims["exp"].(float64)),
			Iss: claims["iss"].(string),
			UID: claims["uid"].(string),
		}, nil
	}

	return nil, ErrInvalidToken
}
