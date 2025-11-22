package handlers

import (
	"database/sql"
	"fmt"
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
	"gofi/internal/services"
	"gofi/internal/types"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type authHandler struct {
	app *app.Application
}

func (h *authHandler) SignUp(c *fiber.Ctx) error {
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
		RoleID:    uuid.Must(uuid.Parse(constant.RoleUser)),
	}

	userVerifyAccount := &models.UserVerifyAccount{}

	err := lib.WithTransaction(h.app.Repositories.User.DB, func(tx *sql.Tx) error {
		err := user.BeforeCreate()
		if err != nil {
			return err
		}

		err = h.app.Repositories.User.InsertExec(tx, user)
		if err != nil {
			return err
		}

		jsonWebToken := jwt.New(&h.app.Config.App)
		token, expiresIn, err := jsonWebToken.Generate(&jwt.JWTPayload{
			UID:       user.ID.String(),
			Secret:    h.app.Config.App.JWTSecret,
			ExpiresAt: "1", // 1 day
		})
		if err != nil {
			return err
		}

		userVerifyAccount.ID = user.ID
		userVerifyAccount.Token = token
		userVerifyAccount.ExpiresAt = time.Unix(expiresIn, 0)

		return h.app.Repositories.UserVerifyAccount.InsertExec(tx, userVerifyAccount)
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	link := fmt.Sprintf("%s/verify?token=%s", h.app.Config.App.ClientURL, userVerifyAccount.Token)

	emailForm := struct {
		Fullname string
		Link     string
	}{
		Fullname: strings.Join([]string{dto.FirstName, *dto.LastName}, " "),
		Link:     link,
	}

	_, err = h.app.Services.Email.SendEmail(services.SendEmailParams{
		Subject:      "Verify your email address",
		To:           dto.Email,
		Data:         emailForm,
		HtmlTemplate: "templates/emails/registration.html",
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Sign up successfully",
	})
}

func (h *authHandler) SignIn(c *fiber.Ctx) error {
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

	user, err := h.app.Repositories.User.GetByEmail(dto.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	hash := argon2.New()
	match, err := hash.Compare(*user.Password, dto.Password)
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

	jsonWebToken := jwt.New(&h.app.Config.App)
	token, expiresIn, err := jsonWebToken.Generate(&jwt.JWTPayload{
		UID:       user.ID.String(),
		Secret:    h.app.Config.App.JWTSecret,
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

	err = h.app.Repositories.Session.Insert(session)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(types.ResponseSingleData[any]{
		Message: "Sign in successfully",
		Data: fiber.Map{
			"uid":          user.ID.String(),
			"email":        user.Email,
			"display_name": strings.Join([]string{user.FirstName, *user.LastName}, " "),
			"is_admin":     user.RoleID.String() == constant.RoleAdmin,
			"access_token": token,
		},
	})
}

func (h *authHandler) VerifySession(c *fiber.Ctx) error {
	uid, err := lib.ContextGetUID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := h.app.Repositories.User.Get(uid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(types.ResponseSingleData[*models.User]{
		Message: "Verify session successfully",
		Data:    user,
	})
}

func (h *authHandler) SignOut(c *fiber.Ctx) error {
	jwt := jwt.New(&h.app.Config.App)

	extractToken, err := jwt.ExtractToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": fmt.Sprintf("Unauthorized, %s", err.Error()),
		})
	}

	uid, err := lib.ContextGetUID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.app.Repositories.Session.Delete(uid, extractToken)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Sign out successfully",
	})
}
