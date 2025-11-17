package lib

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ContextGetUID(c *fiber.Ctx) uuid.UUID {
	uid := uuid.MustParse(c.Locals("uid").(string))
	return uid
}

func ContextSetUID(c *fiber.Ctx, uid uuid.UUID) {
	c.Locals("uid", uid.String())
}
