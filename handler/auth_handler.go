package handler

import (
	"fmt"
	"goarif-api/lib"
	"goarif-api/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) RegisterRoutes(route fiber.Router) {
	auth := route.Group("/auth")
	auth.Post("/sign-up", h.SignUp)
	auth.Post("/sign-in", h.SignIn)
	auth.Get("/verify-session", h.VerifySession)
	auth.Post("/sign-out", h.SignOut)
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var req service.AuthSignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "Invalid request body")
	}

	fmt.Println(req)

	err := h.service.SignUp(req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been added", nil)
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var req service.AuthSignInRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "Invalid request body")
	}

	req.IpAddress = c.IP()
	req.UserAgent = c.Get("User-Agent")

	data, err := h.service.SignIn(req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been added", data)
}

func (h *AuthHandler) VerifySession(c *fiber.Ctx) error {
	data, err := h.service.VerifySession(c)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been added", data)
}

func (h *AuthHandler) SignOut(c *fiber.Ctx) error {
	err := h.service.SignOut(c)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been added", nil)
}
