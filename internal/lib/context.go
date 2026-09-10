package lib

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ContextGetUID(c *fiber.Ctx) (uuid.UUID, error) {
	value := c.Locals("uid")
	if value == nil {
		return uuid.Nil, errors.New("can't find get context auth, please check your authorization")
	}

	raw, ok := value.(string)
	if !ok {
		return uuid.Nil, errors.New("invalid uid in context, please check your authorization")
	}

	uid, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, errors.New("invalid uid in context, please check your authorization")
	}

	return uid, nil
}

func ContextSetUID(c *fiber.Ctx, uid uuid.UUID) {
	c.Locals("uid", uid.String())
}

func ContextParamUUID(c *fiber.Ctx, key string) (uuid.UUID, error) {
	str := c.Params(key)
	return uuid.Parse(str)
}
