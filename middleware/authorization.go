package middleware

import (
	"database/sql"
	"gofi/config"
	"gofi/database"
	"gofi/database/model"
	"gofi/lib"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Authorization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		secretKey := config.Env("JWT_SECRET_KEY", "secret")

		db := database.GetDB()
		var session model.Session

		extractToken, err := lib.ExtractToken(c)
		if err != nil {
			return lib.SendUnauthorizedResponse(c, err.Error())
		}

		query := `
			SELECT * FROM "session" WHERE token = $1
		`
		err = db.Get(&session, query, extractToken)
		if err != sql.ErrNoRows && err != nil {
			return lib.SendUnauthorizedResponse(c, err.Error())
		}

		// check session from header token
		if session.ID != uuid.Nil {
			claims, err := lib.VerifyToken(c, secretKey)
			if err != nil {
				return lib.SendUnauthorizedResponse(c, err.Error())
			}

			uid, err := uuid.Parse(claims.UID)
			if err != nil {
				return lib.SendUnauthorizedResponse(c, "invalid UID in token")
			}

			if claims.Exp < time.Now().Unix() {
				return lib.SendUnauthorizedResponse(c, "token is invalid")
			}

			c.Locals("uid", uid.String())
		}

		return c.Next()
	}
}
