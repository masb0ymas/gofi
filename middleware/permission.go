package middleware

import (
	"database/sql"
	"gofi/database"
	"gofi/database/model"
	"gofi/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/masb0ymas/go-utils/pkg"
)

func PermissionAccess(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		uid := uuid.MustParse(c.Locals("uid").(string))

		var user model.User

		query := `
			SELECT * FROM "user" WHERE id = $1 AND deleted_at IS NULL
		`
		err := db.Get(&user, query, uid)
		if err != sql.ErrNoRows && err != nil {
			return lib.SendUnauthorizedResponse(c, err.Error())
		}

		errType := "permitted access error:"
		errMessage := "you are not allowed"

		if user.ID != uuid.Nil && !pkg.Contains(roles, user.RoleID.String()) {
			return lib.SendUnauthorizedResponse(c, errType+" "+errMessage)
		}

		return c.Next()
	}
}

func NotPermittedAccess(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := database.GetDB()
		uid := uuid.MustParse(c.Locals("uid").(string))

		var user model.User

		query := `
			SELECT * FROM "user" WHERE id = $1 AND deleted_at IS NULL
		`
		err := db.Get(&user, query, uid)
		if err != sql.ErrNoRows && err != nil {
			return lib.SendUnauthorizedResponse(c, err.Error())
		}

		errType := "not permitted access error:"
		errMessage := "you are not allowed"

		if user.ID != uuid.Nil && pkg.Contains(roles, user.RoleID.String()) {
			return lib.SendUnauthorizedResponse(c, errType+" "+errMessage)
		}

		return c.Next()
	}
}
