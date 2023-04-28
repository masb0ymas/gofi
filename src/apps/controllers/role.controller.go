package controllers

import (
	"gofi/src/apps/services"
	"gofi/src/pkg/config"
	"gofi/src/pkg/modules"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func FindAllRole(c *fiber.Ctx) error {
	db := config.GetDB()

	roleService := services.NewRoleService(db)
	data, total, err := roleService.FindAll(c)

	if err != nil {
		panic(err)
	}

	httpResponse := modules.HttpResponse(modules.Response{
		Data:  data,
		Total: total,
	})

	return c.Status(http.StatusOK).JSON(httpResponse)
}
