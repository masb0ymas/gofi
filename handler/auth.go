package handler

import (
	"fmt"
	"gofi/config"
	"gofi/lib"
	"gofi/middleware"
	"gofi/service"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	new_route := route.Group("/auth")
	new_route.Post("/sign-up", h.SignUp)
	new_route.Post("/sign-in", h.SignIn)
	new_route.Get("/verify", h.VerifyToken)
	new_route.Get("/verify-session", middleware.Authorization(), h.VerifySession)
	new_route.Post("/sign-out", middleware.Authorization(), h.SignOut)
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var req service.AuthSignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "invalid request body")
	}

	uid, token, err := h.service.SignUp(req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	serverURL := config.Env("APP_SERVER_URL", "http://localhost:8000")

	// Send Email
	data := struct {
		Fullname string
		Link     string
	}{
		Fullname: req.Fullname,
		Link:     fmt.Sprintf("%s/v1/auth/verify?uid=%s&token-verify=%s", serverURL, uid, url.QueryEscape(*token)),
	}

	err = lib.SendEmail(lib.SendEmailParams{
		Subject:          fmt.Sprintf("Registration from %s", req.Fullname),
		To:               req.Email,
		Data:             data,
		FilenameTemplate: "registration.html",
	})
	if err != nil {
		return lib.SendBadRequestResponse(c, err.Error())
	}

	return lib.SendSuccessResponse(c, "data has been added", nil)
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var req service.AuthSignInRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "invalid request body")
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

func (h *AuthHandler) VerifyToken(c *fiber.Ctx) error {
	UID := c.Query("uid", "")
	Token := c.Query("token-verify", "")

	clientURL := config.Env("APP_CLIENT_URL", "http://localhost:3000")
	failureURL := fmt.Sprintf("%s/verify/failure", clientURL)
	successURL := fmt.Sprintf("%s/verify/success", clientURL)

	// Validate UID
	uid, err := uuid.Parse(UID)
	if err != nil {
		return c.Redirect(failureURL, 303) // Status See Other
	}

	// Verify token
	err = h.service.VerifyToken(uid, Token)
	if err != nil {
		return c.Redirect(failureURL, 303) // Status See Other
	}

	// Success redirect
	return c.Redirect(successURL, 303) // Status See Other
}

func (h *AuthHandler) SignOut(c *fiber.Ctx) error {
	err := h.service.SignOut(c)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been added", nil)
}
