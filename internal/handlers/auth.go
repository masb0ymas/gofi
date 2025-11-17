package handlers

import (
	"net/http"
	"strings"
	"time"

	"gofi/internal/app"
	"gofi/internal/dto"
	"gofi/internal/lib"
	"gofi/internal/lib/argon2"
	"gofi/internal/lib/constant"
	"gofi/internal/lib/jwt"
	"gofi/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type authHandler struct {
	app *app.Application
}

func (r *authHandler) SignUp(c *fiber.Ctx) error {
	var dto dto.AuthSignUp

	if err := lib.ValidateRequestBody(c, &dto); err != nil {
		switch e := err.(type) {
		case *lib.ErrValidationFailed:
			return c.Status(http.StatusBadRequest).JSON(lib.WrapValidationError(e.MessageRecord))
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	user := &models.User{
		Base: models.Base{
			ID: uuid.Must(uuid.NewV7()),
		},
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Phone:     dto.Phone,
		Password:  &dto.Password,
		ActiveAt:  lib.TimePtr(time.Now()),
		RoleID:    uuid.Must(uuid.Parse(constant.RoleUser)),
	}

	err := r.app.Repositories.User.Insert(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Sign up successfully",
	})
}

func (r *authHandler) SignIn(c *fiber.Ctx) error {
	var dto dto.AuthSignIn

	if err := lib.ValidateRequestBody(c, &dto); err != nil {
		switch e := err.(type) {
		case *lib.ErrValidationFailed:
			return c.Status(http.StatusBadRequest).JSON(lib.WrapValidationError(e.MessageRecord))
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	user, err := r.app.Repositories.User.GetByEmail(dto.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	arg := argon2.New()
	match, err := arg.Compare(*user.Password, dto.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if !match {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	jsonWebToken := jwt.New(&r.app.Config.App)
	token, expiresIn, err := jsonWebToken.Generate(&jwt.JWTPayload{
		UID:       user.ID.String(),
		Secret:    r.app.Config.App.JWTSecret,
		ExpiresAt: "1", // 1 day
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	session := &models.Session{
		Base: models.Base{
			ID: uuid.Must(uuid.NewV7()),
		},
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Unix(expiresIn, 0),
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	}

	err = r.app.Repositories.Session.Insert(session)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Sign in successfully",
		"data": fiber.Map{
			"uid":          user.ID.String(),
			"email":        user.Email,
			"display_name": strings.Join([]string{user.FirstName, *user.LastName}, " "),
			"is_admin":     user.RoleID.String() == constant.RoleAdmin,
			"access_token": token,
		},
	})
}

func (r *authHandler) VerifySession(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

func (r *authHandler) SignOut(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}
