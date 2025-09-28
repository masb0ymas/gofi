package lib

import "github.com/gofiber/fiber/v2"

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Metadata   interface{} `json:"metadata,omitempty"`
}

func SendResponse(c *fiber.Ctx, statusCode int, message string, data interface{}, metadata interface{}) error {
	return c.Status(statusCode).JSON(Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Metadata:   metadata,
	})
}

func SendGetResponse(c *fiber.Ctx, message string, data interface{}, metadata interface{}) error {
	return SendResponse(c, fiber.StatusOK, message, data, metadata)
}

func SendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return SendResponse(c, fiber.StatusOK, message, data, nil)
}

func SendCreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return SendResponse(c, fiber.StatusCreated, message, data, nil)
}

func SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return SendResponse(c, statusCode, message, nil, nil)
}

func SendBadRequestResponse(c *fiber.Ctx, message string) error {
	return SendErrorResponse(c, fiber.StatusBadRequest, message)
}

func SendNotFoundResponse(c *fiber.Ctx, message string) error {
	return SendErrorResponse(c, fiber.StatusNotFound, message)
}

func SendUnauthorizedResponse(c *fiber.Ctx, message string) error {
	return SendErrorResponse(c, fiber.StatusUnauthorized, message)
}

func SendForbiddenResponse(c *fiber.Ctx, message string) error {
	return SendErrorResponse(c, fiber.StatusForbidden, message)
}

func SendInternalServerErrorResponse(c *fiber.Ctx, err error) error {
	return SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
}
