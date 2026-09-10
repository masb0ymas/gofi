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
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, ErrInvalidToken
		}

		iss, ok := claims["iss"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		uid, ok := claims["uid"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		return &JWTClaims{
			Exp: int64(exp),
			Iss: iss,
			UID: uid,
		}, nil
	}

	return nil, ErrInvalidToken
}
