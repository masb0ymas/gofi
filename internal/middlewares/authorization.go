package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"gofi/internal/config"
	"gofi/internal/lib"
	"gofi/internal/lib/apperr"
)

// Authorization requires a live Authula session cookie and stores the
// signed-in user's id for handlers (see lib.ContextGetUID). It mirrors the
// session plugin's "session.auth" check for routes outside the Authula
// handler; sliding renewal still happens on Authula's own routes (e.g. /me).
func (m *Middlewares) Authorization(c fiber.Ctx) error {
	token := c.Cookies(config.SessionCookieName)
	if token == "" {
		return apperr.ErrUnauthorized
	}

	services := m.app.Auth.CoreServices()
	session, err := services.SessionService.GetByToken(c.Context(), services.TokenService.Hash(token))
	if err != nil || session == nil || session.ExpiresAt.Before(time.Now().UTC()) {
		return apperr.ErrUnauthorized
	}

	uid, err := uuid.Parse(session.UserID)
	if err != nil {
		return apperr.ErrUnauthorized
	}
	lib.ContextSetUID(c, uid)

	return c.Next()
}
