package controllers

import (
	"fmt"
	"gofi/src/apps/services"
	"gofi/src/database/entities"
	"gofi/src/pkg/config"
	"gofi/src/pkg/modules"
	"net/http"

	"github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
)

func FindAllUser(c *fiber.Ctx) error {
	db := config.GetDB()
	var newData []entities.UserAllEntity

	userService := services.NewUserService(db)
	data, total, err := userService.FindAll(c)

	if err != nil {
		fmt.Print(err)
		panic(err)
	}

	for _, v := range data {
		var record entities.UserAllEntity

		dto.Map(&record, v)
		newData = append(newData, record)
	}

	httpResponse := modules.HttpResponse(modules.Response{
		Data:  newData,
		Total: total,
	})

	return c.Status(http.StatusOK).JSON(httpResponse)
}
