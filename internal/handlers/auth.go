package handlers

import (
	"database/sql"
	"errors"
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
	"gofi/internal/repositories"
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
		AppName  string
	}{
		Fullname: strings.Join([]string{dto.FirstName, *dto.LastName}, " "),
		Link:     link,
		AppName:  h.app.Config.App.Name,
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

	expiresAt := time.Now().Add(time.Hour * 24 * 60) // 60 days
	rt := lib.NewRefreshToken(&h.app.Config.App)
	refToken := rt.Generate(user.ID.String(), expiresAt.Unix())

	refreshToken := &models.RefreshToken{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    user.ID,
		Token:     refToken,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	err = lib.WithTransaction(h.app.Repositories.Session.DB, func(tx *sql.Tx) error {
		err = h.app.Repositories.Session.InsertExec(tx, session)
		if err != nil {
			return err
		}

		return h.app.Repositories.RefreshToken.InsertExec(tx, refreshToken)
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(types.ResponseSingleData[any]{
		Message: "Sign in successfully",
		Data: fiber.Map{
			"uid":           user.ID.String(),
			"email":         user.Email,
			"display_name":  strings.Join([]string{user.FirstName, *user.LastName}, " "),
			"is_admin":      user.RoleID.String() == constant.RoleAdmin,
			"access_token":  token,
			"refresh_token": refToken,
		},
	})
}

func (h *authHandler) VerifyRegistration(c *fiber.Ctx) error {
	var dto dto.AuthVerifyRegistration

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

	jsonWebToken := jwt.New(&h.app.Config.App)
	claims, err := jsonWebToken.Verify(dto.Token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	userID := uuid.Must(uuid.Parse(claims.UID))

	userVerifyAccount, err := h.app.Repositories.UserVerifyAccount.Get(userID, dto.Token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if userVerifyAccount.ExpiresAt.Before(time.Now()) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Token expired",
		})
	}

	user, err := h.app.Repositories.User.Get(userVerifyAccount.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user.ActiveAt = lib.TimePtr(time.Now())

	err = h.app.Repositories.User.Update(user.ID, user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Verify registration successfully",
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

func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	var dto dto.AuthRefreshToken

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

	rt, err := h.app.Repositories.RefreshToken.Get(uid, dto.Token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if rt.ExpiresAt.Before(time.Now()) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Token expired",
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

	extractToken, err := jsonWebToken.ExtractToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": fmt.Sprintf("Unauthorized, %s", err.Error()),
		})
	}

	session, err := h.app.Repositories.Session.GetByUserToken(uid, extractToken)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	session.Token = token
	session.ExpiresAt = time.Unix(expiresIn, 0)
	session.IPAddress = c.IP()
	session.UserAgent = c.Get("User-Agent")

	err = h.app.Repositories.Session.Update(session.ID, session)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(types.ResponseSingleData[any]{
		Message: "Refresh token successfully",
		Data: fiber.Map{
			"uid":           user.ID.String(),
			"email":         user.Email,
			"display_name":  strings.Join([]string{user.FirstName, *user.LastName}, " "),
			"is_admin":      user.RoleID.String() == constant.RoleAdmin,
			"access_token":  token,
			"refresh_token": dto.Token,
		},
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

func (h *authHandler) GoogleAuthURL(c *fiber.Ctx) error {
	url, err := h.app.Services.Google.AuthCodeURL()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"url": url,
	})
}

func (h *authHandler) GoogleAuthCallback(c *fiber.Ctx) error {
	var dto dto.AuthGoogle

	if err := lib.ValidateRequestQuery(c, &dto); err != nil {
		switch e := err.(type) {
		case *lib.ErrValidationFailed:
			return c.Status(http.StatusBadRequest).JSON(lib.WrapValidationError(e.MessageRecord))
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	rawURL := fmt.Sprintf("%s/v1/auth/google/callback?state=%s&code=%s", h.app.Config.App.ServerURL, dto.State, dto.Code)
	state, code, err := h.app.Services.Google.URLParse(rawURL)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	result, err := h.app.Services.Google.Authenticate(state, code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var user *models.User
	var accessToken string

	user, err = h.app.Repositories.User.GetByEmail(result.UserInfo.Email)
	// if errors
	if err != nil {
		// if users not found, create a new user
		if errors.Is(err, repositories.ErrRecordNotFound) {
			// create user from oauth google
			accessToken, err = h.createUserOAuthGoogle(c, result)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message": err.Error(),
				})
			}
		} else {
			fmt.Printf("IS AN ERRORS, THROW ERROR")
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	// if not errors
	if err == nil {
		fmt.Printf("IS NOT AN ERRORS, JUST UPDATE THE USER OAUTH")
		// if user is exists, just insert session and update user oauth
		accessToken, err = h.updateUserOAuthGoogle(c, user, result)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	return c.Status(http.StatusOK).JSON(types.ResponseSingleData[any]{
		Message: "Google auth successfully",
		Data: fiber.Map{
			"uid":           result.UserInfo.ID,
			"email":         result.UserInfo.Email,
			"display_name":  result.UserInfo.Name,
			"is_admin":      false,
			"access_token":  accessToken,
			"refresh_token": result.Token.RefreshToken,
			"token":         result.Token,
			"user":          result.UserInfo,
		},
	})
}

func (h *authHandler) createUserOAuthGoogle(c *fiber.Ctx, authResponse *services.AuthenticateResponse) (token string, err error) {
	jsonWebToken := jwt.New(&h.app.Config.App)

	user := &models.User{
		Base: models.Base{
			ID:        uuid.Must(uuid.NewV7()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		FirstName: authResponse.UserInfo.GivenName,
		LastName:  lib.StringPtr(authResponse.UserInfo.FamilyName),
		Email:     authResponse.UserInfo.Email,
		ActiveAt:  lib.TimePtr(time.Now()),
		RoleID:    uuid.MustParse(constant.RoleUser),
	}

	err = lib.WithTransaction(h.app.Repositories.User.DB, func(tx *sql.Tx) error {
		err := h.app.Repositories.User.InsertExec(tx, user)
		if err != nil {
			return err
		}

		token, expiresIn, err := jsonWebToken.Generate(&jwt.JWTPayload{
			UID:       user.ID.String(),
			Secret:    h.app.Config.App.JWTSecret,
			ExpiresAt: "1", // 1 day
		})
		if err != nil {
			return err
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

		err = h.app.Repositories.Session.InsertExec(tx, session)
		if err != nil {
			return err
		}

		userOAuth := &models.UserOAuth{
			ID:           uuid.Must(uuid.NewV7()),
			UserID:       user.ID,
			Provider:     "google",
			AccessToken:  authResponse.Token.AccessToken,
			RefreshToken: lib.StringPtr(authResponse.Token.RefreshToken),
			ExpiresAt:    time.Unix(authResponse.Token.Expiry.Unix(), 0),
		}

		return h.app.Repositories.UserOAuth.InsertExec(tx, userOAuth)
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

func (h *authHandler) updateUserOAuthGoogle(c *fiber.Ctx, user *models.User, authResponse *services.AuthenticateResponse) (token string, err error) {
	jsonWebToken := jwt.New(&h.app.Config.App)

	err = lib.WithTransaction(h.app.Repositories.User.DB, func(tx *sql.Tx) error {
		userOAuth, err := h.app.Repositories.UserOAuth.GetByUserProviderExec(tx, user.ID, "google")
		if err != nil {
			return err
		}

		userOAuth.AccessToken = authResponse.Token.AccessToken
		userOAuth.RefreshToken = lib.StringPtr(authResponse.Token.RefreshToken)
		userOAuth.ExpiresAt = time.Unix(authResponse.Token.Expiry.Unix(), 0)

		token, expiresIn, err := jsonWebToken.Generate(&jwt.JWTPayload{
			UID:       user.ID.String(),
			Secret:    h.app.Config.App.JWTSecret,
			ExpiresAt: "1", // 1 day
		})
		if err != nil {
			return err
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

		err = h.app.Repositories.Session.InsertExec(tx, session)
		if err != nil {
			return err
		}

		return h.app.Repositories.UserOAuth.UpdateExec(tx, userOAuth.ID, userOAuth)
	})

	return token, err
}
