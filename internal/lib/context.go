package lib

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ContextGetUID(c *fiber.Ctx) (uuid.UUID, error) {
	if c.Locals("uid") != nil {
		uid := uuid.MustParse(c.Locals("uid").(string))
		return uid, nil
	}

	return uuid.Nil, errors.New("can't find get context auth, please check your authorization")
}

func ContextSetUID(c *fiber.Ctx, uid uuid.UUID) {
	c.Locals("uid", uid.String())
}

func ContextParamUUID(c *fiber.Ctx, key string) (uuid.UUID, error) {
	str := c.Params(key)
	return uuid.Parse(str)
}
