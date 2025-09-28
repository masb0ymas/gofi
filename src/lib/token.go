package lib

import (
	"errors"
	"fmt"
	"gofi/src/config"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	UID       uuid.UUID
	SecretKey string
	ExpiresAt string
}

type JWTClaims struct {
	Exp int64
	Iss string
	UID string
}

func GenerateToken(payload *Payload) (string, int64, error) {
	var APP_NAME = config.Env("APP_NAME", "gofi")

	expiresIn, err := strconv.Atoi(payload.ExpiresAt) // expires in days
	if err != nil {
		log.Fatal(err)
	}

	var EXPIRES_TOKEN = time.Now().Add(time.Hour * 24 * time.Duration(expiresIn)).Unix()

	claims := jwt.MapClaims{
		"iss": APP_NAME,
		"exp": EXPIRES_TOKEN,
		"uid": payload.UID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(payload.SecretKey))
	if err != nil {
		log.Fatal(err)
	}

	return t, EXPIRES_TOKEN, nil
}

func ExtractToken(c *fiber.Ctx) (string, error) {
	getTokenByQuery := c.Query("token")
	if getTokenByQuery != "" {
		return getTokenByQuery, nil
	}

	getTokenByCookie := c.Cookies("token")
	if getTokenByCookie != "" {
		return getTokenByCookie, nil
	}

	getTokenByHeader := c.Get("Authorization")
	if getTokenByHeader != "" {
		// Split the token into parts
		parts := strings.Split(getTokenByHeader, " ")
		if len(parts) != 2 {
			return "", errors.New("invalid token format")
		}

		return parts[1], nil
	}

	return "", nil
}

func VerifyToken(c *fiber.Ctx, secretKey string) (*JWTClaims, error) {
	extractToken, err := ExtractToken(c)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(extractToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		errMsg := fmt.Sprintf("Token verification failed: %v. Please ensure your token is valid and not expired.", err)
		return nil, errors.New(errMsg)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &JWTClaims{
			Exp: int64(claims["exp"].(float64)),
			Iss: claims["iss"].(string),
			UID: claims["uid"].(string),
		}, nil
	}

	return nil, err
}
