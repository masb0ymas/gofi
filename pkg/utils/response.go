package utils

import "github.com/gofiber/fiber/v2"

func FailureResponse(code int32, message string, errors interface{}) interface{} {
	return fiber.Map{
		"code":    code,
		"message": message,
		"errors":  errors,
	}
}

func SuccessResponse(code int32, message string, data interface{}) interface{} {
	return fiber.Map{
		"code":    code,
		"message": message,
		"data":    data,
	}
}
