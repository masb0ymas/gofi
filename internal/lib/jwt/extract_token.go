package jwt

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (j *JWT) ExtractToken(c *fiber.Ctx) (string, error) {
	contextQuery := c.Query("token")
	if contextQuery != "" {
		return contextQuery, nil
	}

	contextCookie := c.Cookies("token")
	if contextCookie != "" {
		return contextCookie, nil
	}

	contextHeader := c.Get("Authorization")
	if contextHeader != "" {
		// Bearer <token>
		parts := strings.Split(contextHeader, " ")
		if len(parts) != 2 {
			return "", errors.New("invalid token format")
		}

		if parts[0] != "Bearer" || parts[1] == "" {
			return "", errors.New("invalid token format")
		}

		return parts[1], nil
	}

	return "", errors.New("token not found")
}
